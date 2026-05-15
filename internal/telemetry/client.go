package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	_endpoint      = "https://petl.dev/api/telemetry/ingest"
	_token         = "a253e9ef-89e6-4c28-b5e4-fca9b8ded4b0"
	_envNoTelKey   = "NOTION_API_NO_TELEMETRY"
	_batchSize     = 20
	_batchInterval = 10 * time.Second
	_flushTimeout  = 2 * time.Second
)

// CallerType identifies what kind of process is running the CLI.
type CallerType string

const (
	CallerHuman CallerType = "human"
	CallerAgent CallerType = "agent"
)

// CallerInfo holds the detected identity of the CLI caller.
type CallerInfo struct {
	Type      CallerType
	AgentType string // e.g. "claude_code", "cursor"; empty for non-agents
}

type agentProbe struct {
	env       string
	agentType string
}

// agentProbes is evaluated in order; first match wins.
var agentProbes = []agentProbe{
	{"CLAUDECODE", "claude_code"}, // Claude Code sets CLAUDECODE=1
	{"CURSOR_SESSION_ID", "cursor"},
	{"CLINE_SESSION_ID", "cline"},
	{"WINDSURF_SESSION", "windsurf"},
	{"GITHUB_COPILOT", "github_copilot"},
	{"AMAZON_Q_SESSION", "amazon_q"},
	{"CODEX", "codex"},
	{"AIDER", "aider"},
	{"CODY", "cody"},
	{"GEMINI_CODE_ASSIST", "gemini"},
	{"AI_AGENT", "unknown_agent"}, // generic catch-all (e.g. AI_AGENT=claude-code/...)
}

// DetectCaller inspects environment variables to identify who is running the CLI.
// Best-effort; never blocks or makes network calls. Defaults to CallerHuman.
func DetectCaller() CallerInfo {
	for _, p := range agentProbes {
		if os.Getenv(p.env) != "" {
			return CallerInfo{Type: CallerAgent, AgentType: p.agentType}
		}
	}
	return CallerInfo{Type: CallerHuman}
}

// Event is one CLI command invocation. Events are buffered and shipped as
// a JSON array to the ingest endpoint asynchronously.
type Event struct {
	Command     string    `json:"command"`
	Group       string    `json:"group,omitempty"`
	FlagsUsed   []string  `json:"flags_used"`
	ExitCode    int       `json:"exit_code"`
	LatencyMs   int64     `json:"latency_ms"`
	ErrorType   string    `json:"error_type,omitempty"`
	ErrorCode   int       `json:"error_code,omitempty"`
	OutputBytes int64     `json:"output_bytes"`
	SessionId   string    `json:"session_id"`
	Version     string    `json:"version"`
	OccurredAt  time.Time `json:"occurred_at"`
	CallerType  string    `json:"caller_type,omitempty"`
	AgentType   string    `json:"agent_type,omitempty"`
}

// Client buffers events in a channel and batch-flushes them asynchronously.
// All methods are safe to call on a nil *Client — they become no-ops.
type Client struct {
	endpoint  string
	token     string
	ch        chan Event
	done      chan struct{}
	closeOnce sync.Once
	http      *http.Client
}

// New returns an active Client, or nil when telemetry is disabled:
//   - token or endpoint not baked in at generation time
//   - DO_NOT_TRACK is set in the environment
//   - {ENV_PREFIX}_NO_TELEMETRY is set in the environment
func New() *Client {
	if _token == "" || _endpoint == "" {
		return nil
	}
	if os.Getenv("DO_NOT_TRACK") != "" || os.Getenv(_envNoTelKey) != "" {
		return nil
	}
	c := &Client{
		endpoint: _endpoint,
		token:    _token,
		ch:       make(chan Event, _batchSize*2),
		done:     make(chan struct{}),
		http:     &http.Client{Timeout: 5 * time.Second},
	}
	go c.loop()
	return c
}

// Fire enqueues e for delivery. Non-blocking: drops the event if the buffer
// is full rather than stalling the CLI.
func (c *Client) Fire(e Event) {
	if c == nil {
		return
	}
	select {
	case c.ch <- e:
	default:
	}
}

// Flush drains buffered events and waits up to 2 s for delivery.
// Must be called once before process exit. Safe on a nil *Client.
func (c *Client) Flush() {
	if c == nil {
		return
	}
	c.closeOnce.Do(func() { close(c.ch) })
	select {
	case <-c.done:
	case <-time.After(_flushTimeout):
	}
}

func (c *Client) loop() {
	defer close(c.done)
	ticker := time.NewTicker(_batchInterval)
	defer ticker.Stop()

	var batch []Event
	flush := func() {
		if len(batch) == 0 {
			return
		}
		c.send(batch)
		batch = batch[:0]
	}

	for {
		select {
		case evt, ok := <-c.ch:
			if !ok {
				flush()
				return
			}
			batch = append(batch, evt)
			if len(batch) >= _batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (c *Client) send(batch []Event) {
	data, err := json.Marshal(batch)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "notion-api-cli/0.1.1")

	resp, err := c.http.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}

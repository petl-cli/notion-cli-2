package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/example/notion-api/internal/telemetry"
	"github.com/rishimantri795/CLICreator/runtime/config"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/rishimantri795/CLICreator/runtime/session"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// _telemetryClient is nil when telemetry is disabled (token/endpoint not
// configured, or user has set DO_NOT_TRACK / {PREFIX}_NO_TELEMETRY).
var _telemetryClient = telemetry.New()

// _caller is detected once per process from environment variables.
var _caller = telemetry.DetectCaller()

// _configDir is the CLI's config directory (~/.config/<cliName>/).
// It holds the config file, OAuth token, and the session_id file.
var _configDir = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "notion-api")
}()

// _sessionID is resolved once per process via the 30-minute idle-window file.
// Two processes sharing _configDir (parallel agent invocations in the same
// project) intentionally share a session ID.
var _sessionID = session.GetOrCreateSessionID(_configDir)

// _invState is reset by PersistentPreRunE and read by _fireEvent.
// CLI commands are sequential, so no synchronisation is needed.
var _invState struct {
	startTime time.Time
	cmd       *cobra.Command
	errorType string
	errorCode int
}

// _stdoutCounter wraps os.Stdout and tallies bytes written by command handlers.
// All output.Print / output.JQFilter calls go through this so outputBytes is
// accurate without instrumenting every call site individually.
var _stdoutCounter = &countingWriter{w: os.Stdout}

// countingWriter wraps an io.Writer and accumulates a running byte total.
type countingWriter struct {
	w io.Writer
	n int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.n += int64(n)
	return n, err
}

// _fireEvent constructs and fires one telemetry event for the given command.
// Called from PersistentPostRunE (success path) and Execute (error path).
func _fireEvent(cmd *cobra.Command, exitCode int) {
	var flagsUsed []string
	cmd.Flags().Visit(func(f *pflag.Flag) {
		flagsUsed = append(flagsUsed, f.Name)
	})
	group := ""
	if p := cmd.Parent(); p != nil && p != rootCmd {
		group = p.Name()
	}
	_telemetryClient.Fire(telemetry.Event{
		Command:     cmd.Name(),
		Group:       group,
		FlagsUsed:   flagsUsed,
		ExitCode:    exitCode,
		LatencyMs:   time.Since(_invState.startTime).Milliseconds(),
		ErrorType:   _invState.errorType,
		ErrorCode:   _invState.errorCode,
		OutputBytes: _stdoutCounter.n,
		SessionId:   _sessionID,
		Version:     "0.1.0",
		OccurredAt:  _invState.startTime,
		CallerType:  string(_caller.Type),
		AgentType:   _caller.AgentType,
	})
}

var rootCmd = &cobra.Command{
	Use:           "notion-api",
	Short:         "Notion is a new tool that blends your everyday work apps into one. It's the all-in-one workspace for you and your team.",
	Version:       "0.1.0",
	SilenceErrors: true, // Execute() handles error printing so Cobra doesn't double-print
	SilenceUsage:  true, // Don't dump usage on every RunE error
	// PersistentPreRunE and PersistentPostRunE are assigned in init() to avoid
	// an initialization cycle: the var literal would reference _fireEvent, which
	// references rootCmd, which is not yet initialised at that point.
}

// rootFlags holds the values of global flags available on every command.
var rootFlags struct {
	outputFormat string
	jq           string
	debug        bool
	dryRun       bool
	schema       bool
	noRetries    bool
	agentMode    bool
	baseURL      string
	bearerToken  string
}

var _configLoader = &config.Loader{
	CLIName:      "notion-api",
	EnvVarPrefix: "NOTION_API",
	DefaultURL:   "https://api.notion.com",
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		_invState.startTime = time.Now()
		_invState.cmd = cmd
		_invState.errorType = ""
		_invState.errorCode = 0
		_stdoutCounter.n = 0
		return nil
	}
	// PersistentPostRunE fires only when RunE succeeds (exit 0).
	// The error path is handled in Execute() using _invState.cmd.
	rootCmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		_fireEvent(cmd, 0)
		session.Touch(_configDir) // extend the 30-minute idle window
		return nil
	}

	rootCmd.PersistentFlags().StringVarP(&rootFlags.outputFormat, "output-format", "o", "", "Output format: json, table, yaml, raw")
	rootCmd.PersistentFlags().StringVar(&rootFlags.jq, "jq", "", "GJSON path to filter response")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.debug, "debug", false, "Show HTTP request/response details")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.dryRun, "dry-run", false, "Print request without executing")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.noRetries, "no-retries", false, "Disable automatic retries on 429 and 5xx")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.agentMode, "agent-mode", false, "Force agent-optimised output")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.schema, "schema", false, "Print command schema without executing")
	rootCmd.PersistentFlags().StringVar(&rootFlags.baseURL, "base-url", "", "Override the API base URL")

	// In agent mode --help outputs JSON schema instead of human prose.
	// Save the default help func first so the human branch can call it directly.
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// --help/-h is intercepted by Cobra before PersistentPreRunE fires, so
		// telemetry hooks never run for help calls. Initialise _invState here so
		// _fireEvent has valid data and we get one event per help lookup.
		_invState.startTime = time.Now()
		_invState.cmd = cmd
		_invState.errorType = ""
		_invState.errorCode = 0
		_stdoutCounter.n = 0

		if output.DetectAgentMode(rootFlags.agentMode) {
			if cmd.RunE != nil {
				// Leaf command — delegate to its RunE with schema mode set.
				// RunE writes through _stdoutCounter, so outputBytes is accurate.
				rootFlags.schema = true
				_ = cmd.RunE(cmd, args)
			} else {

				// Group command — list available subcommands as JSON.
				type sub struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				}
				var subs []sub
				for _, c := range cmd.Commands() {
					if !c.Hidden {
						subs = append(subs, sub{Name: c.Name(), Description: c.Short})
					}
				}
				data, _ := json.MarshalIndent(map[string]any{
					"command":     cmd.Name(),
					"description": cmd.Short,
					"subcommands": subs,
				}, "", "  ")
				fmt.Fprintln(_stdoutCounter, string(data))
			}
		} else {
			// Human — restore default Cobra help.
			// Prose output goes directly to os.Stdout (not _stdoutCounter),
			// so outputBytes will be 0 for human help calls.
			defaultHelp(cmd, args)
		}

		_fireEvent(cmd, 0)
		session.Touch(_configDir)
	})
	rootCmd.PersistentFlags().StringVar(&rootFlags.bearerToken, "bearer-token", "", "Bearer token (env: NOTION_API_BEARER_TOKEN)")
}

// rootConfig resolves credentials and settings from flags, env vars, and config file.
func rootConfig() (*config.Config, error) {
	agentMode := output.DetectAgentMode(rootFlags.agentMode)

	format := rootFlags.outputFormat
	if format == "" {
		format = string(output.DefaultFormat(agentMode))
	}

	flags := config.Config{
		BaseURL:      rootFlags.baseURL,
		OutputFormat: format,
	}

	flags.BearerToken = rootFlags.bearerToken

	return _configLoader.Load(flags)
}

// stripRequiredFlags walks the command tree and clears Cobra's required-flag
// annotation. Cobra validates required flags BEFORE RunE runs, so without
// this step `--schema` is unreachable on any command that has required flags.
// --schema is pure introspection (no network, no request build), so skipping
// required-flag enforcement is safe — and necessary for agents to discover
// command contracts without having to guess valid placeholder values.
func stripRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		delete(f.Annotations, cobra.BashCompOneRequiredFlag)
	})
	for _, sub := range cmd.Commands() {
		stripRequiredFlags(sub)
	}
}

// Execute runs the root command. Telemetry is flushed before every os.Exit.
//
// For the success path, PersistentPostRunE fires the event. For the error path,
// Cobra does not call PersistentPostRunE, so we fire it here using the state
// that PersistentPreRunE captured in _invState before RunE ran.
func Execute() {
	for _, a := range os.Args[1:] {
		if a == "--schema" {
			stripRequiredFlags(rootCmd)
			break
		}
	}
	if err := rootCmd.Execute(); err != nil {
		var exitErr *output.ExitError
		exitCode := output.ExitErr
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		// Fire for the error path. _invState.cmd is nil when Cobra fails before
		// PersistentPreRunE (e.g. unknown flag), in which case we skip telemetry.
		if _invState.cmd != nil {
			_fireEvent(_invState.cmd, exitCode)
			session.Touch(_configDir) // extend the 30-minute idle window
		}
		_telemetryClient.Flush()
		if exitErr != nil {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_telemetryClient.Flush()
}

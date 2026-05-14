package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var commentGetListCmd = &cobra.Command{
	Use:   "get-list",
	Short: "Retrieve comments",
	RunE:  runCommentGetList,
}

var commentGetListFlags struct {
	notionVersion string
	blockId       string
	startCursor   string
	pageSize      int
}

func init() {
	commentGetListCmd.Flags().StringVar(&commentGetListFlags.notionVersion, "notion-version", "", "")
	commentGetListCmd.MarkFlagRequired("notion-version")
	commentGetListCmd.Flags().StringVar(&commentGetListFlags.blockId, "block-id", "", "Identifier for a Notion block or page")
	commentGetListCmd.MarkFlagRequired("block-id")
	commentGetListCmd.Flags().StringVar(&commentGetListFlags.startCursor, "start-cursor", "", "If supplied, this endpoint will return a page of results starting after the cursor provided. If not supplied, this endpoint will return the first page of results.")
	commentGetListCmd.Flags().IntVar(&commentGetListFlags.pageSize, "page-size", 0, "The number of items from the full list desired in the response. Maximum: 100")

	commentCmd.AddCommand(commentGetListCmd)
}

func runCommentGetList(cmd *cobra.Command, args []string) error {
	// --schema: print full input/output type contract without making any network call.
	if rootFlags.schema {
		type flagSchema struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Required    bool   `json:"required"`
			Location    string `json:"location"`
			Description string `json:"description,omitempty"`
		}
		var flags []flagSchema
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    true,
			Location:    "header",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "block-id",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "Identifier for a Notion block or page",
		})
		flags = append(flags, flagSchema{
			Name:        "start-cursor",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "If supplied, this endpoint will return a page of results starting after the cursor provided. If not supplied, this endpoint will return the first page of results.",
		})
		flags = append(flags, flagSchema{
			Name:        "page-size",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "The number of items from the full list desired in the response. Maximum: 100",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "200",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "403",
		})

		schema := map[string]any{
			"command":     "get-list",
			"description": "Retrieve comments",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v1/comments",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         true,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{},
				"impact":       "low",
			},
			"requires_auth": false,
		}
		data, _ := json.MarshalIndent(schema, "", "  ")
		fmt.Fprintln(_stdoutCounter, string(data))
		return nil
	}

	cfg, err := rootConfig()
	if err != nil {
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	client := httpclient.New(cfg.BaseURL, cfg.AuthProvider())
	client.Debug = rootFlags.debug
	client.DryRun = rootFlags.dryRun
	if rootFlags.noRetries {
		client.RetryConfig.MaxRetries = 0
	}

	// Build path params
	pathParams := map[string]string{}

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v1/comments", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("block-id") {
		req.QueryParams["block_id"] = fmt.Sprintf("%v", commentGetListFlags.blockId)
	}
	if cmd.Flags().Changed("start-cursor") {
		req.QueryParams["start_cursor"] = fmt.Sprintf("%v", commentGetListFlags.startCursor)
	}
	if cmd.Flags().Changed("page-size") {
		req.QueryParams["page_size"] = fmt.Sprintf("%v", commentGetListFlags.pageSize)
	}

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", commentGetListFlags.notionVersion)
	}

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			_invState.errorType = "timeout"
		} else {
			_invState.errorType = "network_error"
		}
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode >= 500 {
			_invState.errorType = "http_5xx"
		} else {
			_invState.errorType = "http_4xx"
		}
		_invState.errorCode = resp.StatusCode
		e := output.HTTPError(resp.StatusCode, resp.Body)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if rootFlags.jq != "" {
		return output.JQFilter(_stdoutCounter, resp.Body, rootFlags.jq)
	}
	return output.Print(_stdoutCounter, resp.Body, output.Format(cfg.OutputFormat))
}

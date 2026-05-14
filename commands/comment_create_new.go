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

var commentCreateNewCmd = &cobra.Command{
	Use:   "create-new",
	Short: "Create comment",
	RunE:  runCommentCreateNew,
}

var commentCreateNewFlags struct {
	notionVersion string
	parent        string
	discussionId  string
	richText      string
	body          string
}

func init() {
	commentCreateNewCmd.Flags().StringVar(&commentCreateNewFlags.notionVersion, "notion-version", "", "")
	commentCreateNewCmd.MarkFlagRequired("notion-version")
	commentCreateNewCmd.Flags().StringVar(&commentCreateNewFlags.parent, "parent", "", "A [page parent](/reference/database#page-parent). Either this or a discussion_id is required (not both)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	commentCreateNewCmd.Flags().StringVar(&commentCreateNewFlags.discussionId, "discussion-id", "", "A UUID identifier for a discussion thread. Either this or a parent object is required (not both)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	commentCreateNewCmd.Flags().StringVar(&commentCreateNewFlags.richText, "rich-text", "", "A [rich text object](ref:rich-text)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	commentCreateNewCmd.Flags().StringVar(&commentCreateNewFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	commentCmd.AddCommand(commentCreateNewCmd)
}

func runCommentCreateNew(cmd *cobra.Command, args []string) error {
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
			Name:        "parent",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A [page parent](/reference/database#page-parent). Either this or a discussion_id is required (not both)",
		})
		flags = append(flags, flagSchema{
			Name:        "discussion-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A UUID identifier for a discussion thread. Either this or a parent object is required (not both)",
		})
		flags = append(flags, flagSchema{
			Name:        "rich-text",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "A [rich text object](ref:rich-text)",
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
			"command":     "create-new",
			"description": "Create comment",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v1/comments",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   false,
				"reversible":   true,
				"side_effects": []string{"creates_resource"},
				"impact":       "medium",
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
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v1/comments", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", commentCreateNewFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if commentCreateNewFlags.body != "" {
		if err := json.Unmarshal([]byte(commentCreateNewFlags.body), &bodyMap); err != nil {
			_invState.errorType = "parse_error"
			cliErr := &output.CLIError{
				Error:    true,
				Code:     "validation_error",
				Message:  fmt.Sprintf("invalid JSON in --body: %v", err),
				ExitCode: output.ExitValidation,
			}
			cliErr.Write(os.Stderr)
			return output.NewExitError(cliErr)
		}
	}
	// Individual flags overlay onto body (flags take precedence over --body JSON)
	if cmd.Flags().Changed("parent") {
		bodyMap["parent"] = commentCreateNewFlags.parent
	}
	if cmd.Flags().Changed("discussion-id") {
		bodyMap["discussion_id"] = commentCreateNewFlags.discussionId
	}
	if cmd.Flags().Changed("rich-text") {
		bodyMap["rich_text"] = commentCreateNewFlags.richText
	}
	req.Body = bodyMap

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

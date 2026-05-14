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

var blockAppendChildrenCmd = &cobra.Command{
	Use:   "append-children",
	Short: "Append block children",
	RunE:  runBlockAppendChildren,
}

var blockAppendChildrenFlags struct {
	blockId       string
	notionVersion string
	children      []string
	after         string
	body          string
}

func init() {
	blockAppendChildrenCmd.Flags().StringVar(&blockAppendChildrenFlags.blockId, "block-id", "", "Identifier for a [block](ref:block). Also accepts a [page](ref:page) ID.")
	blockAppendChildrenCmd.MarkFlagRequired("block-id")
	blockAppendChildrenCmd.Flags().StringVar(&blockAppendChildrenFlags.notionVersion, "notion-version", "", "")
	blockAppendChildrenCmd.MarkFlagRequired("notion-version")
	blockAppendChildrenCmd.Flags().StringSliceVar(&blockAppendChildrenFlags.children, "children", nil, "Child content to append to a container block as an array of [block objects](ref:block)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	blockAppendChildrenCmd.Flags().StringVar(&blockAppendChildrenFlags.after, "after", "", "The ID of the existing block that the new block should be appended after.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	blockAppendChildrenCmd.Flags().StringVar(&blockAppendChildrenFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	blockCmd.AddCommand(blockAppendChildrenCmd)
}

func runBlockAppendChildren(cmd *cobra.Command, args []string) error {
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
			Name:        "block-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Identifier for a [block](ref:block). Also accepts a [page](ref:page) ID.",
		})
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    true,
			Location:    "header",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "children",
			Type:        "array",
			Required:    true,
			Location:    "body",
			Description: "Child content to append to a container block as an array of [block objects](ref:block)",
		})
		flags = append(flags, flagSchema{
			Name:        "after",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The ID of the existing block that the new block should be appended after.",
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
			Status:      "400",
			ContentType: "application/json",
			Description: "400",
		})

		schema := map[string]any{
			"command":     "append-children",
			"description": "Append block children",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/v1/blocks/{block_id}/children",
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
				"side_effects": []string{"mutates_resource"},
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
	pathParams["block_id"] = fmt.Sprintf("%v", blockAppendChildrenFlags.blockId)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/v1/blocks/{block_id}/children", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", blockAppendChildrenFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if blockAppendChildrenFlags.body != "" {
		if err := json.Unmarshal([]byte(blockAppendChildrenFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("children") {
		bodyMap["children"] = blockAppendChildrenFlags.children
	}
	if cmd.Flags().Changed("after") {
		bodyMap["after"] = blockAppendChildrenFlags.after
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

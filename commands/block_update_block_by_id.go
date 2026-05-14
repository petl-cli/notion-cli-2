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

var blockUpdateBlockByIdCmd = &cobra.Command{
	Use:   "update-block-by-id",
	Short: "Update a block",
	RunE:  runBlockUpdateBlockById,
}

var blockUpdateBlockByIdFlags struct {
	blockId       string
	notionVersion string
	archived      bool
	body          string
}

func init() {
	blockUpdateBlockByIdCmd.Flags().StringVar(&blockUpdateBlockByIdFlags.blockId, "block-id", "", "Identifier for a Notion block")
	blockUpdateBlockByIdCmd.MarkFlagRequired("block-id")
	blockUpdateBlockByIdCmd.Flags().StringVar(&blockUpdateBlockByIdFlags.notionVersion, "notion-version", "", "")
	blockUpdateBlockByIdCmd.Flags().BoolVar(&blockUpdateBlockByIdFlags.archived, "archived", false, "Set to true to archive (delete) a block. Set to false to un-archive (restore) a block.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	blockUpdateBlockByIdCmd.Flags().StringVar(&blockUpdateBlockByIdFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	blockCmd.AddCommand(blockUpdateBlockByIdCmd)
}

func runBlockUpdateBlockById(cmd *cobra.Command, args []string) error {
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
			Description: "Identifier for a Notion block",
		})
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    false,
			Location:    "header",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "type",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "The [block object `type`](ref:block#block-object-keys) value with the properties to be updated. Currently only `text` (for supported block types) and `checked` (for `to_do` blocks) fields can be updated.",
		})
		flags = append(flags, flagSchema{
			Name:        "archived",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set to true to archive (delete) a block. Set to false to un-archive (restore) a block.",
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
			"command":     "update-block-by-id",
			"description": "Update a block",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/v1/blocks/{block_id}",
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
	pathParams["block_id"] = fmt.Sprintf("%v", blockUpdateBlockByIdFlags.blockId)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/v1/blocks/{block_id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", blockUpdateBlockByIdFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if blockUpdateBlockByIdFlags.body != "" {
		if err := json.Unmarshal([]byte(blockUpdateBlockByIdFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("archived") {
		bodyMap["archived"] = blockUpdateBlockByIdFlags.archived
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

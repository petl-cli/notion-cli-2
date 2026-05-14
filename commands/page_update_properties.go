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

var pageUpdatePropertiesCmd = &cobra.Command{
	Use:   "update-properties",
	Short: "Update page properties",
	RunE:  runPageUpdateProperties,
}

var pageUpdatePropertiesFlags struct {
	pageId        string
	notionVersion string
	properties    string
	archived      bool
	icon          string
	cover         string
	body          string
}

func init() {
	pageUpdatePropertiesCmd.Flags().StringVar(&pageUpdatePropertiesFlags.pageId, "page-id", "", "The identifier for the Notion page to be updated.")
	pageUpdatePropertiesCmd.MarkFlagRequired("page-id")
	pageUpdatePropertiesCmd.Flags().StringVar(&pageUpdatePropertiesFlags.notionVersion, "notion-version", "", "")
	pageUpdatePropertiesCmd.Flags().StringVar(&pageUpdatePropertiesFlags.properties, "properties", "", "The property values to update for the page. The keys are the names or IDs of the property and the values are property values. If a page property ID is not included, then it is not changed.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageUpdatePropertiesCmd.Flags().BoolVar(&pageUpdatePropertiesFlags.archived, "archived", false, "Whether the page is archived (deleted). Set to true to archive a page. Set to false to un-archive (restore) a page.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageUpdatePropertiesCmd.Flags().StringVar(&pageUpdatePropertiesFlags.icon, "icon", "", "A page icon for the page. Supported types are [external file object](https://developers.notion.com/reference/file-object) or [emoji object](https://developers.notion.com/reference/emoji-object).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageUpdatePropertiesCmd.Flags().StringVar(&pageUpdatePropertiesFlags.cover, "cover", "", "A cover image for the page. Only [external file objects](https://developers.notion.com/reference/file-object) are supported.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageUpdatePropertiesCmd.Flags().StringVar(&pageUpdatePropertiesFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	pageCmd.AddCommand(pageUpdatePropertiesCmd)
}

func runPageUpdateProperties(cmd *cobra.Command, args []string) error {
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
			Name:        "page-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "The identifier for the Notion page to be updated.",
		})
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    false,
			Location:    "header",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "properties",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The property values to update for the page. The keys are the names or IDs of the property and the values are property values. If a page property ID is not included, then it is not changed.",
		})
		flags = append(flags, flagSchema{
			Name:        "archived",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Whether the page is archived (deleted). Set to true to archive a page. Set to false to un-archive (restore) a page.",
		})
		flags = append(flags, flagSchema{
			Name:        "icon",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A page icon for the page. Supported types are [external file object](https://developers.notion.com/reference/file-object) or [emoji object](https://developers.notion.com/reference/emoji-object).",
		})
		flags = append(flags, flagSchema{
			Name:        "cover",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A cover image for the page. Only [external file objects](https://developers.notion.com/reference/file-object) are supported.",
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
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "404",
		})
		responses = append(responses, responseSchema{
			Status:      "429",
			ContentType: "application/json",
			Description: "429",
		})

		schema := map[string]any{
			"command":     "update-properties",
			"description": "Update page properties",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/v1/pages/{page_id}",
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
	pathParams["page_id"] = fmt.Sprintf("%v", pageUpdatePropertiesFlags.pageId)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/v1/pages/{page_id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", pageUpdatePropertiesFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if pageUpdatePropertiesFlags.body != "" {
		if err := json.Unmarshal([]byte(pageUpdatePropertiesFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("properties") {
		bodyMap["properties"] = pageUpdatePropertiesFlags.properties
	}
	if cmd.Flags().Changed("archived") {
		bodyMap["archived"] = pageUpdatePropertiesFlags.archived
	}
	if cmd.Flags().Changed("icon") {
		bodyMap["icon"] = pageUpdatePropertiesFlags.icon
	}
	if cmd.Flags().Changed("cover") {
		bodyMap["cover"] = pageUpdatePropertiesFlags.cover
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

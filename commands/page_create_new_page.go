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

var pageCreateNewPageCmd = &cobra.Command{
	Use:   "create-new-page",
	Short: "Create a page",
	RunE:  runPageCreateNewPage,
}

var pageCreateNewPageFlags struct {
	notionVersion string
	parent        string
	properties    string
	children      []string
	icon          string
	cover         string
	body          string
}

func init() {
	pageCreateNewPageCmd.Flags().StringVar(&pageCreateNewPageFlags.notionVersion, "notion-version", "", "")
	pageCreateNewPageCmd.MarkFlagRequired("notion-version")
	pageCreateNewPageCmd.Flags().StringVar(&pageCreateNewPageFlags.parent, "parent", "", "The parent page or database where the new page is inserted, represented as a JSON object with a `page_id` or `database_id` key, and the corresponding ID.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageCreateNewPageCmd.Flags().StringVar(&pageCreateNewPageFlags.properties, "properties", "", "The values of the page’s properties. If the `parent` is a database, then the schema must match the parent database’s properties. If the `parent` is a page, then the only valid object key is `title`.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageCreateNewPageCmd.Flags().StringSliceVar(&pageCreateNewPageFlags.children, "children", nil, "The content to be rendered on the new page, represented as an array of [block objects](https://developers.notion.com/reference/block).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageCreateNewPageCmd.Flags().StringVar(&pageCreateNewPageFlags.icon, "icon", "", "The icon of the new page. Either an [emoji object](https://developers.notion.com/reference/emoji-object) or an [external file object](https://developers.notion.com/reference/file-object)..")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageCreateNewPageCmd.Flags().StringVar(&pageCreateNewPageFlags.cover, "cover", "", "The cover image of the new page, represented as a [file object](https://developers.notion.com/reference/file-object).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	pageCreateNewPageCmd.Flags().StringVar(&pageCreateNewPageFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	pageCmd.AddCommand(pageCreateNewPageCmd)
}

func runPageCreateNewPage(cmd *cobra.Command, args []string) error {
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
			Required:    true,
			Location:    "body",
			Description: "The parent page or database where the new page is inserted, represented as a JSON object with a `page_id` or `database_id` key, and the corresponding ID.",
		})
		flags = append(flags, flagSchema{
			Name:        "properties",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "The values of the page’s properties. If the `parent` is a database, then the schema must match the parent database’s properties. If the `parent` is a page, then the only valid object key is `title`.",
		})
		flags = append(flags, flagSchema{
			Name:        "children",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "The content to be rendered on the new page, represented as an array of [block objects](https://developers.notion.com/reference/block).",
		})
		flags = append(flags, flagSchema{
			Name:        "icon",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The icon of the new page. Either an [emoji object](https://developers.notion.com/reference/emoji-object) or an [external file object](https://developers.notion.com/reference/file-object)..",
		})
		flags = append(flags, flagSchema{
			Name:        "cover",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The cover image of the new page, represented as a [file object](https://developers.notion.com/reference/file-object).",
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
			"command":     "create-new-page",
			"description": "Create a page",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v1/pages",
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
		Path:        httpclient.SubstitutePath("/v1/pages", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", pageCreateNewPageFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if pageCreateNewPageFlags.body != "" {
		if err := json.Unmarshal([]byte(pageCreateNewPageFlags.body), &bodyMap); err != nil {
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
		bodyMap["parent"] = pageCreateNewPageFlags.parent
	}
	if cmd.Flags().Changed("properties") {
		bodyMap["properties"] = pageCreateNewPageFlags.properties
	}
	if cmd.Flags().Changed("children") {
		bodyMap["children"] = pageCreateNewPageFlags.children
	}
	if cmd.Flags().Changed("icon") {
		bodyMap["icon"] = pageCreateNewPageFlags.icon
	}
	if cmd.Flags().Changed("cover") {
		bodyMap["cover"] = pageCreateNewPageFlags.cover
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

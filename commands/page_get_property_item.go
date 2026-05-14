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

var pageGetPropertyItemCmd = &cobra.Command{
	Use:   "get-property-item",
	Short: "Retrieve a page property item",
	RunE:  runPageGetPropertyItem,
}

var pageGetPropertyItemFlags struct {
	pageId        string
	propertyId    string
	pageSize      int
	startCursor   string
	notionVersion string
}

func init() {
	pageGetPropertyItemCmd.Flags().StringVar(&pageGetPropertyItemFlags.pageId, "page-id", "", "Identifier for a Notion page")
	pageGetPropertyItemCmd.MarkFlagRequired("page-id")
	pageGetPropertyItemCmd.Flags().StringVar(&pageGetPropertyItemFlags.propertyId, "property-id", "", "Identifier for a page [property](https://developers.notion.com/reference/page#all-property-values)")
	pageGetPropertyItemCmd.MarkFlagRequired("property-id")
	pageGetPropertyItemCmd.Flags().IntVar(&pageGetPropertyItemFlags.pageSize, "page-size", 0, "For paginated properties. The max number of property item objects on a page. The default size is 100")
	pageGetPropertyItemCmd.Flags().StringVar(&pageGetPropertyItemFlags.startCursor, "start-cursor", "", "For paginated properties.")
	pageGetPropertyItemCmd.Flags().StringVar(&pageGetPropertyItemFlags.notionVersion, "notion-version", "", "")

	pageCmd.AddCommand(pageGetPropertyItemCmd)
}

func runPageGetPropertyItem(cmd *cobra.Command, args []string) error {
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
			Description: "Identifier for a Notion page",
		})
		flags = append(flags, flagSchema{
			Name:        "property-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Identifier for a page [property](https://developers.notion.com/reference/page#all-property-values)",
		})
		flags = append(flags, flagSchema{
			Name:        "page-size",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "For paginated properties. The max number of property item objects on a page. The default size is 100",
		})
		flags = append(flags, flagSchema{
			Name:        "start-cursor",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "For paginated properties.",
		})
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    false,
			Location:    "header",
			Description: "",
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

		schema := map[string]any{
			"command":     "get-property-item",
			"description": "Retrieve a page property item",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v1/pages/{page_id}/properties/{property_id}",
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
	pathParams["page_id"] = fmt.Sprintf("%v", pageGetPropertyItemFlags.pageId)
	pathParams["property_id"] = fmt.Sprintf("%v", pageGetPropertyItemFlags.propertyId)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v1/pages/{page_id}/properties/{property_id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("page-size") {
		req.QueryParams["page_size"] = fmt.Sprintf("%v", pageGetPropertyItemFlags.pageSize)
	}
	if cmd.Flags().Changed("start-cursor") {
		req.QueryParams["start_cursor"] = fmt.Sprintf("%v", pageGetPropertyItemFlags.startCursor)
	}

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", pageGetPropertyItemFlags.notionVersion)
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

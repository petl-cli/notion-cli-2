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

var databaseExecuteQueryCmd = &cobra.Command{
	Use:   "execute-query",
	Short: "Query a database",
	RunE:  runDatabaseExecuteQuery,
}

var databaseExecuteQueryFlags struct {
	databaseId       string
	filterProperties string
	notionVersion    string
	filter           string
	sorts            []string
	startCursor      string
	pageSize         int
	body             string
}

func init() {
	databaseExecuteQueryCmd.Flags().StringVar(&databaseExecuteQueryFlags.databaseId, "database-id", "", "Identifier for a Notion database.")
	databaseExecuteQueryCmd.MarkFlagRequired("database-id")
	databaseExecuteQueryCmd.Flags().StringVar(&databaseExecuteQueryFlags.filterProperties, "filter-properties", "", "A list of page property value IDs associated with the database. Use this param to limit the response to a specific page property value or values for pages that meet the `filter` criteria.")
	databaseExecuteQueryCmd.Flags().StringVar(&databaseExecuteQueryFlags.notionVersion, "notion-version", "", "")
	databaseExecuteQueryCmd.MarkFlagRequired("notion-version")
	databaseExecuteQueryCmd.Flags().StringVar(&databaseExecuteQueryFlags.filter, "filter", "", "When supplied, limits which pages are returned based on the [filter conditions](ref:post-database-query-filter).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseExecuteQueryCmd.Flags().StringSliceVar(&databaseExecuteQueryFlags.sorts, "sorts", nil, "When supplied, orders the results based on the provided [sort criteria](ref:post-database-query-sort).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseExecuteQueryCmd.Flags().StringVar(&databaseExecuteQueryFlags.startCursor, "start-cursor", "", "When supplied, returns a page of results starting after the cursor provided. If not supplied, this endpoint will return the first page of results.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseExecuteQueryCmd.Flags().IntVar(&databaseExecuteQueryFlags.pageSize, "page-size", 0, "The number of items from the full list desired in the response. Maximum: 100")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseExecuteQueryCmd.Flags().StringVar(&databaseExecuteQueryFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	databaseCmd.AddCommand(databaseExecuteQueryCmd)
}

func runDatabaseExecuteQuery(cmd *cobra.Command, args []string) error {
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
			Name:        "database-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Identifier for a Notion database.",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-properties",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "A list of page property value IDs associated with the database. Use this param to limit the response to a specific page property value or values for pages that meet the `filter` criteria.",
		})
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    true,
			Location:    "header",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "filter",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "When supplied, limits which pages are returned based on the [filter conditions](ref:post-database-query-filter).",
		})
		flags = append(flags, flagSchema{
			Name:        "sorts",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "When supplied, orders the results based on the provided [sort criteria](ref:post-database-query-sort).",
		})
		flags = append(flags, flagSchema{
			Name:        "start-cursor",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "When supplied, returns a page of results starting after the cursor provided. If not supplied, this endpoint will return the first page of results.",
		})
		flags = append(flags, flagSchema{
			Name:        "page-size",
			Type:        "integer",
			Required:    false,
			Location:    "body",
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
			Status:      "400",
			ContentType: "application/json",
			Description: "400",
		})

		schema := map[string]any{
			"command":     "execute-query",
			"description": "Query a database",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v1/databases/{database_id}/query",
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
	pathParams["database_id"] = fmt.Sprintf("%v", databaseExecuteQueryFlags.databaseId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v1/databases/{database_id}/query", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("filter-properties") {
		req.QueryParams["filter_properties"] = fmt.Sprintf("%v", databaseExecuteQueryFlags.filterProperties)
	}

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", databaseExecuteQueryFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if databaseExecuteQueryFlags.body != "" {
		if err := json.Unmarshal([]byte(databaseExecuteQueryFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("filter") {
		bodyMap["filter"] = databaseExecuteQueryFlags.filter
	}
	if cmd.Flags().Changed("sorts") {
		bodyMap["sorts"] = databaseExecuteQueryFlags.sorts
	}
	if cmd.Flags().Changed("start-cursor") {
		bodyMap["start_cursor"] = databaseExecuteQueryFlags.startCursor
	}
	if cmd.Flags().Changed("page-size") {
		bodyMap["page_size"] = databaseExecuteQueryFlags.pageSize
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

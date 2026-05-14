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

var searchByTitleCmd = &cobra.Command{
	Use:   "by-title",
	Short: "Search by title",
	RunE:  runSearchByTitle,
}

var searchByTitleFlags struct {
	notionVersion string
	query         string
	startCursor   string
	pageSize      int
	body          string
}

func init() {
	searchByTitleCmd.Flags().StringVar(&searchByTitleFlags.notionVersion, "notion-version", "", "")
	searchByTitleCmd.MarkFlagRequired("notion-version")
	searchByTitleCmd.Flags().StringVar(&searchByTitleFlags.query, "query", "", "The text that the API compares page and database titles against.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	searchByTitleCmd.Flags().StringVar(&searchByTitleFlags.startCursor, "start-cursor", "", "A `cursor` value returned in a previous response that If supplied, limits the response to results starting after the `cursor`. If not supplied, then the first page of results is returned. Refer to [pagination](https://developers.notion.com/reference/intro#pagination) for more details.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	searchByTitleCmd.Flags().IntVar(&searchByTitleFlags.pageSize, "page-size", 0, "The number of items from the full list to include in the response. Maximum: `100`.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	searchByTitleCmd.Flags().StringVar(&searchByTitleFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	searchCmd.AddCommand(searchByTitleCmd)
}

func runSearchByTitle(cmd *cobra.Command, args []string) error {
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
			Name:        "query",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The text that the API compares page and database titles against.",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "A set of criteria, `direction` and `timestamp` keys, that orders the results. The **only** supported timestamp value is `\"last_edited_time\"`. Supported `direction` values are `\"ascending\"` and `\"descending\"`. If `sort` is not provided, then the most recently edited results are returned first.",
		})
		flags = append(flags, flagSchema{
			Name:        "filter",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "A set of criteria, `value` and `property` keys, that limits the results to either only pages or only databases. Possible `value` values are `\"page\"` or `\"database\"`. The only supported `property` value is `\"object\"`.",
		})
		flags = append(flags, flagSchema{
			Name:        "start-cursor",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A `cursor` value returned in a previous response that If supplied, limits the response to results starting after the `cursor`. If not supplied, then the first page of results is returned. Refer to [pagination](https://developers.notion.com/reference/intro#pagination) for more details.",
		})
		flags = append(flags, flagSchema{
			Name:        "page-size",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "The number of items from the full list to include in the response. Maximum: `100`.",
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
			Status:      "429",
			ContentType: "application/json",
			Description: "429",
		})

		schema := map[string]any{
			"command":     "by-title",
			"description": "Search by title",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v1/search",
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
		Path:        httpclient.SubstitutePath("/v1/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", searchByTitleFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if searchByTitleFlags.body != "" {
		if err := json.Unmarshal([]byte(searchByTitleFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("query") {
		bodyMap["query"] = searchByTitleFlags.query
	}
	if cmd.Flags().Changed("start-cursor") {
		bodyMap["start_cursor"] = searchByTitleFlags.startCursor
	}
	if cmd.Flags().Changed("page-size") {
		bodyMap["page_size"] = searchByTitleFlags.pageSize
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

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

var databaseCreateNewDatabaseCmd = &cobra.Command{
	Use:   "create-new-database",
	Short: "Create a database",
	RunE:  runDatabaseCreateNewDatabase,
}

var databaseCreateNewDatabaseFlags struct {
	notionVersion string
	title         []string
	parent        string
	properties    string
	body          string
}

func init() {
	databaseCreateNewDatabaseCmd.Flags().StringVar(&databaseCreateNewDatabaseFlags.notionVersion, "notion-version", "", "")
	databaseCreateNewDatabaseCmd.MarkFlagRequired("notion-version")
	databaseCreateNewDatabaseCmd.Flags().StringSliceVar(&databaseCreateNewDatabaseFlags.title, "title", nil, "Title of database as it appears in Notion. An array of [rich text objects](ref:rich-text).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseCreateNewDatabaseCmd.Flags().StringVar(&databaseCreateNewDatabaseFlags.parent, "parent", "", "A [page parent](/reference/database#page-parent)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseCreateNewDatabaseCmd.Flags().StringVar(&databaseCreateNewDatabaseFlags.properties, "properties", "", "Property schema of database. The keys are the names of properties as they appear in Notion and the values are [property schema objects](https://developers.notion.com/reference/property-schema-object).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseCreateNewDatabaseCmd.Flags().StringVar(&databaseCreateNewDatabaseFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	databaseCmd.AddCommand(databaseCreateNewDatabaseCmd)
}

func runDatabaseCreateNewDatabase(cmd *cobra.Command, args []string) error {
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
			Name:        "title",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Title of database as it appears in Notion. An array of [rich text objects](ref:rich-text).",
		})
		flags = append(flags, flagSchema{
			Name:        "parent",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "A [page parent](/reference/database#page-parent)",
		})
		flags = append(flags, flagSchema{
			Name:        "properties",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Property schema of database. The keys are the names of properties as they appear in Notion and the values are [property schema objects](https://developers.notion.com/reference/property-schema-object).",
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
			"command":     "create-new-database",
			"description": "Create a database",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v1/databases",
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
		Path:        httpclient.SubstitutePath("/v1/databases", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", databaseCreateNewDatabaseFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if databaseCreateNewDatabaseFlags.body != "" {
		if err := json.Unmarshal([]byte(databaseCreateNewDatabaseFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("title") {
		bodyMap["title"] = databaseCreateNewDatabaseFlags.title
	}
	if cmd.Flags().Changed("parent") {
		bodyMap["parent"] = databaseCreateNewDatabaseFlags.parent
	}
	if cmd.Flags().Changed("properties") {
		bodyMap["properties"] = databaseCreateNewDatabaseFlags.properties
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

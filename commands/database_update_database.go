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

var databaseUpdateDatabaseCmd = &cobra.Command{
	Use:   "update-database",
	Short: "Update a database",
	RunE:  runDatabaseUpdateDatabase,
}

var databaseUpdateDatabaseFlags struct {
	databaseId    string
	notionVersion string
	title         []string
	description   []string
	properties    string
	body          string
}

func init() {
	databaseUpdateDatabaseCmd.Flags().StringVar(&databaseUpdateDatabaseFlags.databaseId, "database-id", "", "identifier for a Notion database")
	databaseUpdateDatabaseCmd.MarkFlagRequired("database-id")
	databaseUpdateDatabaseCmd.Flags().StringVar(&databaseUpdateDatabaseFlags.notionVersion, "notion-version", "", "")
	databaseUpdateDatabaseCmd.Flags().StringSliceVar(&databaseUpdateDatabaseFlags.title, "title", nil, "An array of [rich text objects](https://developers.notion.com/reference/rich-text) that represents the title of the database that is displayed in the Notion UI. If omitted, then the database title remains unchanged.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseUpdateDatabaseCmd.Flags().StringSliceVar(&databaseUpdateDatabaseFlags.description, "description", nil, "An array of [rich text objects](https://developers.notion.com/reference/rich-text) that represents the description of the database that is displayed in the Notion UI. If omitted, then the database description remains unchanged.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseUpdateDatabaseCmd.Flags().StringVar(&databaseUpdateDatabaseFlags.properties, "properties", "", "The properties of a database to be changed in the request, in the form of a JSON object. If updating an existing property, then the keys are the names or IDs of the properties as they appear in Notion, and the values are [property schema objects](ref:property-schema-object). If adding a new property, then the key is the name of the new database property and the value is a [property schema object](ref:property-schema-object).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	databaseUpdateDatabaseCmd.Flags().StringVar(&databaseUpdateDatabaseFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	databaseCmd.AddCommand(databaseUpdateDatabaseCmd)
}

func runDatabaseUpdateDatabase(cmd *cobra.Command, args []string) error {
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
			Description: "identifier for a Notion database",
		})
		flags = append(flags, flagSchema{
			Name:        "notion-version",
			Type:        "string",
			Required:    false,
			Location:    "header",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "title",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "An array of [rich text objects](https://developers.notion.com/reference/rich-text) that represents the title of the database that is displayed in the Notion UI. If omitted, then the database title remains unchanged.",
		})
		flags = append(flags, flagSchema{
			Name:        "description",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "An array of [rich text objects](https://developers.notion.com/reference/rich-text) that represents the description of the database that is displayed in the Notion UI. If omitted, then the database description remains unchanged.",
		})
		flags = append(flags, flagSchema{
			Name:        "properties",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The properties of a database to be changed in the request, in the form of a JSON object. If updating an existing property, then the keys are the names or IDs of the properties as they appear in Notion, and the values are [property schema objects](ref:property-schema-object). If adding a new property, then the key is the name of the new database property and the value is a [property schema object](ref:property-schema-object).",
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
			"command":     "update-database",
			"description": "Update a database",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/v1/databases/{database_id}",
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
	pathParams["database_id"] = fmt.Sprintf("%v", databaseUpdateDatabaseFlags.databaseId)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/v1/databases/{database_id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("notion-version") {
		req.Headers["Notion-Version"] = fmt.Sprintf("%v", databaseUpdateDatabaseFlags.notionVersion)
	}

	// Request body
	bodyMap := map[string]any{}
	if databaseUpdateDatabaseFlags.body != "" {
		if err := json.Unmarshal([]byte(databaseUpdateDatabaseFlags.body), &bodyMap); err != nil {
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
		bodyMap["title"] = databaseUpdateDatabaseFlags.title
	}
	if cmd.Flags().Changed("description") {
		bodyMap["description"] = databaseUpdateDatabaseFlags.description
	}
	if cmd.Flags().Changed("properties") {
		bodyMap["properties"] = databaseUpdateDatabaseFlags.properties
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

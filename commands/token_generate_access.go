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

var tokenGenerateAccessCmd = &cobra.Command{
	Use:   "generate-access",
	Short: "Create a token",
	RunE:  runTokenGenerateAccess,
}

var tokenGenerateAccessFlags struct {
	code        string
	grantType   string
	redirectUri string
	body        string
}

func init() {
	tokenGenerateAccessCmd.Flags().StringVar(&tokenGenerateAccessFlags.code, "code", "", "A unique random code that Notion generates to authenticate with your service, generated when a user initiates the OAuth flow.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tokenGenerateAccessCmd.Flags().StringVar(&tokenGenerateAccessFlags.grantType, "grant-type", "", "A constant string: \"authorization_code\".")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tokenGenerateAccessCmd.Flags().StringVar(&tokenGenerateAccessFlags.redirectUri, "redirect-uri", "", "The `\"redirect_uri\"` that was provided in the OAuth Domain & URI section of the integration's Authorization settings. Do not include this field if a `\"redirect_uri\"` query param was not included in the Authorization URL provided to users. In most cases, this field is required.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tokenGenerateAccessCmd.Flags().StringVar(&tokenGenerateAccessFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	tokenCmd.AddCommand(tokenGenerateAccessCmd)
}

func runTokenGenerateAccess(cmd *cobra.Command, args []string) error {
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
			Name:        "code",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "A unique random code that Notion generates to authenticate with your service, generated when a user initiates the OAuth flow.",
		})
		flags = append(flags, flagSchema{
			Name:        "grant-type",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "A constant string: \"authorization_code\".",
		})
		flags = append(flags, flagSchema{
			Name:        "redirect-uri",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "The `\"redirect_uri\"` that was provided in the OAuth Domain & URI section of the integration's Authorization settings. Do not include this field if a `\"redirect_uri\"` query param was not included in the Authorization URL provided to users. In most cases, this field is required.",
		})
		flags = append(flags, flagSchema{
			Name:        "external-account",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Required if and only when building [Link Preview](https://developers.notion.com/docs/link-previews) integrations (otherwise ignored). An object with `key` and `name` properties. `key` should be a unique identifier for the account. Notion uses the `key` to determine whether or not the user is re-connecting the same account. `name` should be some way for the user to know which account they used to authenticate with your service. If a user has authenticated Notion with your integration before and `key` is the same but `name` is different, then Notion updates the `name` associated with your integration.",
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
			"command":     "generate-access",
			"description": "Create a token",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v1/oauth/token",
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
			"requires_auth": true,
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
		Path:        httpclient.SubstitutePath("/v1/oauth/token", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if tokenGenerateAccessFlags.body != "" {
		if err := json.Unmarshal([]byte(tokenGenerateAccessFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("code") {
		bodyMap["code"] = tokenGenerateAccessFlags.code
	}
	if cmd.Flags().Changed("grant-type") {
		bodyMap["grant_type"] = tokenGenerateAccessFlags.grantType
	}
	if cmd.Flags().Changed("redirect-uri") {
		bodyMap["redirect_uri"] = tokenGenerateAccessFlags.redirectUri
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

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agentInstructionsCmd prints the llms.txt content at runtime so agents can
// include it in their system prompt without needing a separate file:
//
//	INSTRUCTIONS=$(notion-api agent-instructions)
var agentInstructionsCmd = &cobra.Command{
	Use:   "agent-instructions",
	Short: "Print machine-readable instructions for AI agents (llms.txt format)",
	Long: `Prints a complete description of this CLI's commands, flags, exit codes,
and usage patterns optimised for inclusion in an AI agent's system prompt.

Example:
  # Include in Claude Code context:
  notion-api agent-instructions > CLAUDE.md

  # Capture inline:
  INSTRUCTIONS=$(notion-api agent-instructions)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(agentInstructionsContent)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentInstructionsCmd)
}

// agentInstructionsContent is the full llms.txt baked into the binary at build time.
// Regenerate by re-running the CLI generator against the same OpenAPI spec.
const agentInstructionsContent = `# notion-api

Notion is a new tool that blends your everyday work apps into one. It's the all-in-one workspace for you and your team.

This file is the agent-facing overview of the ` + "`" + `notion-api` + "`" + ` CLI. It explains what the tool does and how to use it well — *not* every flag of every command. For per-command details run ` + "`" + `notion-api <command> --help` + "`" + ` or ` + "`" + `notion-api <command> --schema` + "`" + ` (returns JSON).

## Install

The binary is ` + "`" + `notion-api` + "`" + `. Build from source or download a release.

## Authentication

- Bearer token — set ` + "`" + `NOTION_API_BEARER_TOKEN` + "`" + ` or pass ` + "`" + `--bearer-token <token>` + "`" + `

Run ` + "`" + `notion-api configure` + "`" + ` to see how to set credentials.

## Commands

notion-api groups operations by resource. One or two examples per group below — for the full set of commands and flags use ` + "`" + `notion-api <group> --help` + "`" + `.

### ` + "`" + `block` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api block append-children --block-id <string> --notion-version <string>    # Append block children
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api block get-children --block-id <string> --notion-version <string>    # Retrieve block children
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api block --help` + "`" + `

### ` + "`" + `comment` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api comment create-new --notion-version <string>    # Create comment
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api comment get-list --notion-version <string> --block-id <string>    # Retrieve comments
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api comment --help` + "`" + `

### ` + "`" + `database` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api database create-new-database --notion-version <string>    # Create a database
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api database execute-query --database-id <string> --notion-version <string>    # Query a database
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api database --help` + "`" + `

### ` + "`" + `page` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api page create-new-page --notion-version <string>    # Create a page
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api page get-page --page-id <string> --notion-version <string>    # Retrieve a page
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api page --help` + "`" + `

### ` + "`" + `search` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api search by-title --notion-version <string>    # Search by title
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api search --help` + "`" + `

### ` + "`" + `token` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api token generate-access    # Create a token
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api token --help` + "`" + `

### ` + "`" + `user` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api user get-token-bot-user --notion-version <string>    # Retrieve your token's bot user
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
notion-api user get-user-by-id --user-id <string> --notion-version <string>    # Retrieve a user
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `notion-api user --help` + "`" + `


## Output and parsing

All output is JSON by default and goes to stdout. Errors go to stderr as a single-line JSON object — stdout stays clean for piping.

` + "`" + `` + "`" + `` + "`" + `bash
notion-api <cmd> --jq <path>          # extract fields without jq installed (GJSON syntax)
notion-api <cmd> -o yaml              # change format: json (default), yaml, table, compact, raw, pretty
notion-api <cmd> --dry-run            # print the HTTP request without sending it
notion-api <cmd> --schema             # JSON schema of inputs and outputs for that command
` + "`" + `` + "`" + `` + "`" + `

GJSON path examples:
- ` + "`" + `--jq id` + "`" + ` — scalar
- ` + "`" + `--jq items.#.id` + "`" + ` — every id from an array
- ` + "`" + `--jq "items.#(active==true)#"` + "`" + ` — filter array by condition

## Exit codes

Branch on ` + "`" + `$?` + "`" + ` rather than parsing stderr:

| Code | Meaning |
|------|---------|
| 0 | success |
| 1 | unknown error |
| 2 | auth failed (401 / 403) |
| 3 | not found (404) |
| 4 | validation error (400 / 422) |
| 5 | rate limited (429) |
| 6 | server error (5xx) |
| 7 | network error |

Error JSON shape on stderr:
` + "`" + `` + "`" + `` + "`" + `json
{"error":true,"code":"not_found","status":404,"message":"...","exit_code":3}
` + "`" + `` + "`" + `` + "`" + `

## Common workflows

### List then fetch one
` + "`" + `` + "`" + `` + "`" + `bash
ID=$(notion-api <group> list --jq "items.0.id")
notion-api <group> get --id "$ID"
` + "`" + `` + "`" + `` + "`" + `

### Capture a created ID
` + "`" + `` + "`" + `` + "`" + `bash
ID=$(notion-api <group> create [flags] --jq id)
` + "`" + `` + "`" + `` + "`" + `

### Safe destructive call
` + "`" + `` + "`" + `` + "`" + `bash
notion-api <group> delete --id X --dry-run    # inspect first
notion-api <group> delete --id X              # exit 0 = deleted, 3 = not found
` + "`" + `` + "`" + `` + "`" + `

### Branch on exit code
` + "`" + `` + "`" + `` + "`" + `bash
notion-api <cmd> [flags]
case $? in
  0) : success ;;
  2) echo "fix credentials" ;;
  3) echo "not found, skip" ;;
  5) sleep 10 ; retry ;;
esac
` + "`" + `` + "`" + `` + "`" + `

## Discovering more

` + "`" + `` + "`" + `` + "`" + `bash
notion-api --help                              # top-level groups
notion-api <group> --help                      # commands in a group
notion-api <group> <cmd> --help                # flags + description
notion-api <group> <cmd> --schema              # JSON schema
notion-api agent-instructions                  # this file, embedded in the binary
` + "`" + `` + "`" + `` + "`" + `

`

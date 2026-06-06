# gsuite-mcp

> Model Context Protocol server in Go exposing Google Workspace APIs (Gmail, Calendar, Contacts) with OAuth 2.0, MCP prompts and resources, retry logic, and stdio transport.

**Family:** agent-libs · **Type:** tool · **Lifecycle:** production · **Owner:** harperreed

## What it does
gsuite-mcp is a Go MCP server that provides programmatic access to Google Workspace APIs: Gmail (list/get/send/draft/labels/delete), Calendar (view/create/update/delete/quick-add), and Contacts (list/search/get/create/update/delete). It exposes 19 MCP tools, 8 workflow prompts, and dynamic MCP resources, authenticating via OAuth 2.0 with exponential-backoff retry and a mock "ish mode" for testing without real credentials. The Go module is `github.com/harper/gsuite-mcp`.

## How it fits
- Depends on: — (no internal prime-radiant-inc deps in go.mod)
- Used by: — (consumed by MCP clients over stdio)
- External: Google Workspace APIs (Gmail, Calendar, People) via OAuth 2.0

## Runtime & data
- Runs: MCP server (stdio transport)
- Data in: MCP tool calls, OAuth 2.0 credentials
- Data out: Gmail/Calendar/Contacts reads and mutations

<!-- Maintained by the maintaining-project-map skill. Do not hand-edit; regenerated. -->

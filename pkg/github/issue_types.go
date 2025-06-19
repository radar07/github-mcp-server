package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListIssueTypes(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_issue_types",
			mcp.WithDescription(t("TOOL_LIST_ISSUE_TYPES_DESCRIPTION", "")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_LIST_ISSUE_TYPES_USER_TYPE", "List Issue Types"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("org_name",
				mcp.Required(),
				mcp.Description("The organization name to list issue types for"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			orgName, err := RequiredParam[string](request, "org_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			issueTypes, resp, err := client.Organizations.ListIssueTypes(ctx, orgName)
			if err != nil {
				return nil, fmt.Errorf("failed to get issue types: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list issue types: %s", string(body))), nil
			}

			r, err := json.Marshal(issueTypes)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal issue: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func CreateIssueType() {}

func UpdateIssueType() {}

func DeleteIssueType() {}

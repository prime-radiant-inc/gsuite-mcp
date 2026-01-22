// ABOUTME: Tests for gmail_manage_labels MCP tool handler
// ABOUTME: Validates label management operations (list, get, create, update, delete)

package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleGmailManageLabels_Actions(t *testing.T) {
	t.Setenv("ISH_MODE", "true")

	srv, err := NewServer(context.Background())
	require.NoError(t, err)

	tests := []struct {
		name        string
		args        map[string]interface{}
		expectError bool
		checkResult func(t *testing.T, result map[string]interface{})
		description string
	}{
		{
			name:        "missing action",
			args:        map[string]interface{}{},
			expectError: true,
			description: "should fail when action is missing",
		},
		{
			name: "invalid action",
			args: map[string]interface{}{
				"action": "invalid",
			},
			expectError: true,
			description: "should fail for unknown action",
		},
		{
			name: "list action",
			args: map[string]interface{}{
				"action": "list",
			},
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "list", result["action"])
				assert.NotNil(t, result["labels"])
				assert.NotNil(t, result["count"])
			},
			description: "should return labels list",
		},
		{
			name: "get action without label_id",
			args: map[string]interface{}{
				"action": "get",
			},
			expectError: true,
			description: "should fail when label_id is missing for get",
		},
		{
			name: "get action with label_id",
			args: map[string]interface{}{
				"action":   "get",
				"label_id": "INBOX",
			},
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "get", result["action"])
				assert.NotNil(t, result["label"])
			},
			description: "should return label details",
		},
		{
			name: "create action without name",
			args: map[string]interface{}{
				"action": "create",
			},
			expectError: true,
			description: "should fail when name is missing for create",
		},
		{
			name: "create action with name",
			args: map[string]interface{}{
				"action": "create",
				"name":   "Test/Label",
			},
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "create", result["action"])
				assert.NotNil(t, result["label"])
				assert.NotEmpty(t, result["message"])
			},
			description: "should create label successfully",
		},
		{
			name: "create action with visibility options",
			args: map[string]interface{}{
				"action":                    "create",
				"name":                      "Test/Visible",
				"label_list_visibility":    "labelShow",
				"message_list_visibility":  "show",
			},
			expectError: false,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "create", result["action"])
				assert.NotNil(t, result["label"])
			},
			description: "should create label with visibility options",
		},
		{
			name: "update action without label_id",
			args: map[string]interface{}{
				"action": "update",
				"name":   "New Name",
			},
			expectError: true,
			description: "should fail when label_id is missing for update",
		},
		{
			name: "update action without any field to update",
			args: map[string]interface{}{
				"action":   "update",
				"label_id": "Label_123",
			},
			expectError: true,
			description: "should fail when no update fields provided",
		},
		{
			name: "delete action without label_id",
			args: map[string]interface{}{
				"action": "delete",
			},
			expectError: true,
			description: "should fail when label_id is missing for delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createMockRequest("gmail_manage_labels", tt.args)
			result, err := srv.handleGmailManageLabels(context.Background(), request)

			require.NoError(t, err, "handler should not return Go error")
			assert.NotNil(t, result)

			if tt.expectError {
				assert.True(t, result.IsError, tt.description)
			} else {
				assert.False(t, result.IsError, tt.description)
				if tt.checkResult != nil && len(result.Content) > 0 {
					// Parse JSON result
					if textContent, ok := result.Content[0].(map[string]interface{}); ok {
						if text, ok := textContent["text"].(string); ok {
							var parsed map[string]interface{}
							if err := json.Unmarshal([]byte(text), &parsed); err == nil {
								tt.checkResult(t, parsed)
							}
						}
					}
				}
			}
		})
	}
}

func TestHandleGmailManageLabels_ErrorMessages(t *testing.T) {
	t.Setenv("ISH_MODE", "true")

	srv, err := NewServer(context.Background())
	require.NoError(t, err)

	tests := []struct {
		name           string
		args           map[string]interface{}
		expectedSubstr string
		description    string
	}{
		{
			name:           "missing action provides guidance",
			args:           map[string]interface{}{},
			expectedSubstr: "list, get, create, update, delete",
			description:    "error should list valid actions",
		},
		{
			name: "missing label_id for get suggests list",
			args: map[string]interface{}{
				"action": "get",
			},
			expectedSubstr: "action: list",
			description:    "error should suggest using list action",
		},
		{
			name: "create without name explains nested labels",
			args: map[string]interface{}{
				"action": "create",
			},
			expectedSubstr: "slashes",
			description:    "error should explain nested label syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createMockRequest("gmail_manage_labels", tt.args)
			result, err := srv.handleGmailManageLabels(context.Background(), request)

			require.NoError(t, err)
			assert.True(t, result.IsError)

			// Check that error message contains helpful guidance
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(map[string]interface{}); ok {
					if text, ok := textContent["text"].(string); ok {
						assert.Contains(t, text, tt.expectedSubstr, tt.description)
					}
				}
			}
		})
	}
}

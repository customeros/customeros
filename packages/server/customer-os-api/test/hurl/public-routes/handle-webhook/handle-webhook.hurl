# Test: Handle Webhook
POST {{cos_url}}/flows/v1/hooks
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "integration": "fathom"
}

HTTP 201
[Captures]
hook_url: jsonpath "$.hook.url"
integration: jsonpath "$.hook.integration"
secret: jsonpath "$.hook.secret"
tenant_id: regex "/flows/v1/([^/]+)"
integration_id: regex "/i/([a-zA-Z0-9]+)"
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.hook.url" matches "/flows/v1/[a-zA-Z0-9-]+/i/[a-zA-Z0-9-]+"
jsonpath "$.hook.integration" == "fathom"
jsonpath "$.hook.secret" matches "^[a-zA-Z0-9_-]{32,}$"



POST {{cos_url}}/flows/v1/{{tenant_id}}/i/{{integration_id}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
  "ai_summary": {
    "html_formatted": "<h2>Meeting Purpose</h2>\n<p>Discuss challenges and solutions</p>",
    "markdown_formatted": "## Meeting Purpose\nDiscuss challenges and solutions",
    "plaintext_formatted": "Meeting Purpose\nDiscuss challenges and solutions",
    "sections": "[{'depth': 2, 'title': 'Meeting Purpose', 'html_formatted': '<p>Discuss challenges</p>', 'markdown_formatted': 'Discuss challenges', 'plaintext_formatted': 'Discuss challenges', 'slack_formatted': '> Discuss challenges'}]",
    "slack_formatted": "*Meeting*\n> Discuss challenges",
    "template_name": "General"
  },
  "fathom_user": {
    "email": "test@example.com",
    "name": "Test User"
  },
  "id": "12345",
  "meeting": {
    "external_domains": "[{'domain_name': 'example.com'}, {'domain_name': 'test.com'}]",
    "has_external_invitees": "True",
    "invitees": "[{'name': 'John Doe', 'email': 'john@example.com', 'is_external': false}, {'name': 'Jane Smith', 'email': 'jane@test.com', 'is_external': true}]",
    "join_url": "https://meet.google.com/test",
    "scheduled_duration_in_minutes": "30",
    "scheduled_end_time": "2024-05-10T10:00:00Z",
    "scheduled_start_time": "2024-05-10T09:30:00Z",
    "title": "Test Meeting"
  },
  "recording": {
    "duration_in_minutes": "25.5",
    "share_url": "https://example.com/share",
    "url": "https://example.com/recording"
  }
}
HTTP 403
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Forbidden"

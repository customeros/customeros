# Test: Create Organization
POST {{cos_url}}/customerbase/v1/organizations
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "customId": "{{custom_id}}",
    "icpFit": true,
    "leadSource": "Web Search",
    "linkedinUrl": "https://linkedin.com/in/{{custom_id}}",
    "name": "Hurl-{{custom_id}}",
    "relationship": "CUSTOMER",
    "website": "https://hurl-{{custom_id}}.ai"
}

HTTP/2 201
[Captures]
org_id: jsonpath "$.organization.id"
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.organization.id" exists

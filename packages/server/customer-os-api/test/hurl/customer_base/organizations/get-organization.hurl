# Test: Get Organization
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
    "website": "https://hurl{{custom_id}}.ai"
}

HTTP/2 201
[Captures]
org_id: jsonpath "$.organization.id"
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.organization.id" exists

GET {{cos_url}}/customerbase/v1/organizations/{{org_id}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.organization.id" == "{{org_id}}"
jsonpath "$.organization.customId" == "{{custom_id}}"
jsonpath "$.organization.icpFit" == true
jsonpath "$.organization.leadSource" == "Web Search"
jsonpath "$.organization.linkedinURL" == "https://linkedin.com/in/{{custom_id}}"
jsonpath "$.organization.name" == "Hurl-{{custom_id}}"
jsonpath "$.organization.relationship" == "CUSTOMER"
jsonpath "$.organization.stage" == "ONBOARDING"
jsonpath "$.organization.website" == "https://hurl{{custom_id}}.ai"
jsonpath "$.organization.cosId" exists

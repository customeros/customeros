# Test: Update Primary Link
POST https://api.customeros.ai/customerbase/v1/organizations
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "customId": "{{custom_id}}",
    "icpFit": true,
    "leadSource": "Web Search",
    "linkedinUrl": "https://linkedin.com/in/{{custom_id}}",
    "name": "Hurl-{{custom_id}}",
    "relationship": "CUSTOMER",
    "website": "https://testify{{custom_id}}.ai"
}

HTTP/2 201
[Captures]
org_id: jsonpath "$.organization.id"
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.organization.id" exists

GET https://api.customeros.ai/customerbase/v1/organizations/{{org_id}}
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 200
[Asserts]
jsonpath "$.status" == "success"

PUT https://api.customeros.ai/customerbase/v1/organizations/{{org_id}}/links/stripe/primary
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "externalId": "stripe-1234"
}

HTTP/2 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.organization.externalSystem" == "stripe"
jsonpath "$.organization.externalId" == "stripe-1234"
jsonpath "$.organization.primary" == true
jsonpath "$.organization.organizationId" exists

PUT https://api.customeros.ai/customerbase/v1/organizations/{{org_id}}/links/stripe/primary
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "externalId": "stripe-5678"
}

HTTP/2 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.organization.externalSystem" == "stripe"
jsonpath "$.organization.externalId" == "stripe-5678"
jsonpath "$.organization.primary" == true
jsonpath "$.organization.organizationId" exists

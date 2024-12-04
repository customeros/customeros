# Test: Create Contact with Invalid Email
POST https://api.customeros.ai/customerbase/v1/contacts
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "email": "@example.com"
}

HTTP/2 400
[Asserts]
jsonpath "$.status" exists
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "@example.com | invalid email format"

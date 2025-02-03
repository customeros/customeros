# Test: Create Contact with Role Account
POST {{cos_url}}/customerbase/v1/contacts
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "email": "test@example.com"
}

HTTP/2 *
[Asserts]
status == 400
jsonpath "$.status" exists
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "test@example.com | email is a role account"

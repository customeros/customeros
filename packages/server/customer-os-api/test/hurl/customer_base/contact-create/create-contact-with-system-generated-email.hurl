# Test: Create Contact with System Generated Email
POST {{cos_url}}/customerbase/v1/contacts
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "email": "{{custom_id}}@example.com"
}

HTTP/2 400
[Asserts]
jsonpath "$.status" exists
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "{{custom_id}}@example.com | email is system generated"

# Test: Create Contacts with Invalid Emails
POST {{cos_url}}/customerbase/v1/contacts/bulk
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[
    {
        "email": "test@example.com"
    },
    {
        "email": "@yahoo.com"
    },
    {
        "email": ""
    },
    {
        "email": "{{custom_id}}@example.com"
    }
]

HTTP/2 400
[Asserts]
jsonpath "$.status" exists
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "No valid contacts found in request"

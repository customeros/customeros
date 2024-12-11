# Test: Create Contacts with Mixed Emails
POST {{cos_url}}/customerbase/v1/contacts/bulk
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[
    {
        "email": "hurl-{{random_str}}a@gmail.com"
    },
    {
        "email": "hurl-{{random_str}}b@gmail.com"
    },
    {
        "email": "hurl-{{random_str}}c@gmail.com"
    },
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

HTTP/2 201
[Asserts]
jsonpath "$.status" exists
jsonpath "$.status" == "partial success"
jsonpath "$.requestId" exists
jsonpath "$.summary.total" == 7
jsonpath "$.summary.success" == 3
jsonpath "$.summary.failed" == 4
jsonpath "$.details" count == 4
jsonpath "$.details[0].value" == "test@example.com"
jsonpath "$.details[0].description" == "email is a role account"
jsonpath "$.details[1].value" == "@yahoo.com"
jsonpath "$.details[1].description" == "invalid email format"
jsonpath "$.details[2].value" == ""
jsonpath "$.details[2].description" == "must provide either email or LinkedIn URL"
jsonpath "$.details[3].value" == "{{custom_id}}@example.com"
jsonpath "$.details[3].description" == "email is system generated"

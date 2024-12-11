# Test: Create Contacts with Valid Emails
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
        }
    ]

HTTP/2 201
[Asserts]
jsonpath "$.status" exists
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.summary.total" == 3
jsonpath "$.summary.success" == 3
jsonpath "$.summary.failed" == 0
jsonpath "$.details.value" == ""
jsonpath "$.details.description" == ""

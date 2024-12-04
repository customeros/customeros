# Test: Create Contacts with Valid Emails
POST https://api.customeros.ai/customerbase/v1/contacts/bulk
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "contacts": [
        {
            "email": "hurl-{{random_str}}-one@gmail.com"
        },
        {
            "email": "hurl-{{random_str}}-two@gmail.com"
        },
        {
            "email": "hurl-{{random_str}}-three@gmail.com"
        }
    ]
}

HTTP/2 400
[Asserts]
body == ""
jsonpath "$.status" exists
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.summary.total" == 3
jsonpath "$.summary.success" == 3
jsonpath "$.summary.failed" == 0
jsonpath "$.details.value" == ""
jsonpath "$.details.description" == ""

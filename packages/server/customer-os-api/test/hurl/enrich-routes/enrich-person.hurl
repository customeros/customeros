# Test: Create Contact with Valid Email
POST https://api.customeros.ai/enrich/v1/person
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "email": "hurl-{{random_str}}@gmail.com"
}

HTTP/2 200
[Asserts]
jsonpath "$.status" exists
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.contact.email" == "hurl-{{random_str}}@gmail.com"
jsonpath "$.contact.linkedinUrl" == ""

# Test: Create Contact with Valid LinkedInUrl
POST {{cos_url}}/customerbase/v1/contacts
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "linkedinUrl": "https://www.linkedin.com/in/hurl-{{random_str}}"
}

HTTP/2 *
[Asserts]
status == 200
jsonpath "$.status" exists
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.contact.email" == ""
jsonpath "$.contact.linkedinUrl" == "https://www.linkedin.com/in/hurl-{{random_str}}"

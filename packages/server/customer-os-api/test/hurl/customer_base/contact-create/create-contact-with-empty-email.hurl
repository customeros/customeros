# Test: Create Contact with Empty Email
POST {{cos_url}}/customerbase/v1/contacts
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "email": ""
}

HTTP/2 200
[Asserts]
body == ""

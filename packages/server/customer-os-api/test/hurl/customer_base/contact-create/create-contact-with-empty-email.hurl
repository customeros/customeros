# Test: Create Contact with Empty Email
POST https://api.customeros.ai/customerbase/v1/contacts
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
{
    "email": ""
}

HTTP/2 200
[Asserts]
body == ""

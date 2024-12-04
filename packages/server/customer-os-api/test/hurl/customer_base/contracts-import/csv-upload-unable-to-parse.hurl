# Test unable to parse CSV file at upload
POST https://api.customeros.ai/customerbase/v1/contacts/import
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: contacts.csv; @ contacts-success-unable-to-parse.csv

HTTP 400
[Asserts]
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Unable to parse file"

# Test: Unsuccessful Import Contacts From CSV
POST https://api.customeros.ai/customerbase/v1/contacts/import
Content-Type: multipart/form-data; boundary=boundary
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: contacts.csv; @ contacts-success-unable-to-parse.csv

HTTP 400
[Asserts]
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Unable to parse file"

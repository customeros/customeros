# Test: Unsuccessful Import Contacts From CSV - Missing Headers
POST https://api.customeros.ai/customerbase/v1/contacts/import
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,unsuccessful-import-contacts-from-csv-missing-headers.csv;

HTTP 400
[Asserts]
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Missing required headers: email, linkedin_url"

# Test: Unsuccessful Import Contacts From CSV unable to read file
POST {{cos_url}}/customerbase/v1/contacts/import
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,unsuccessful-import-contacts-from-csv-no-contacts-found.csv;

HTTP 400
[Asserts]
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "No valid contacts found in request"

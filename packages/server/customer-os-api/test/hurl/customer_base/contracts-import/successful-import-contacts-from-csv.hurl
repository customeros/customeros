# Test: Successful Import Contacts From CSV
POST {{cos_url}}/customerbase/v1/contacts/import
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,successful-import-contacts-from-csv.csv;

HTTP 201
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.summary.total" == 2
jsonpath "$.summary.success" == 2
jsonpath "$.summary.failed" == 0
jsonpath "$.details.value" == ""
jsonpath "$.details.description" == ""

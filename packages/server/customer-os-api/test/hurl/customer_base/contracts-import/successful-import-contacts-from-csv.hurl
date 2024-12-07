# Test: Successful Import Contacts From CSV
POST https://api.customeros.ai/customerbase/v1/contacts/import
Content-Type: multipart/form-data; boundary=boundary
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file:contacts.csv;filename=contacts.csv; @ contacts-success.csv

HTTP 201
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.summary.total" == 2
jsonpath "$.summary.success" == 2
jsonpath "$.summary.failed" == 0

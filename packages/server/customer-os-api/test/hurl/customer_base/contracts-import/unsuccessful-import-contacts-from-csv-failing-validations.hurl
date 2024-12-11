# Test: Unsuccessful Import Contacts From CSV - Unable to Read Record
POST {{cos_url}}/customerbase/v1/contacts/import
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,unsuccessful-import-contacts-from-csv-failing-validations.csv;

HTTP 201
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.summary.total" == 4
jsonpath "$.summary.success" == 1
jsonpath "$.summary.failed" == 3
jsonpath "$.details[0].value" == " "
jsonpath "$.details[0].description" == "invalid LinkedIn URL format"
jsonpath "$.details[1].value" == " "
jsonpath "$.details[1].description" == "invalid email format"
jsonpath "$.details[2].value" == "[johndoethree]"
jsonpath "$.details[2].description" == "Unable to read record"

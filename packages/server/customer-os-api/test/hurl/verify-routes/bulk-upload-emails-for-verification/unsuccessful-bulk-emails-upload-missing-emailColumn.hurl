# Test: Unsuccessful bulk emails upload missing emailColumn
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,unsuccessful-bulk-emails-upload-missing-emailColumn.csv;

HTTP *
[Asserts]
status == 400
header "Content-Type" contains "application/json"
jsonpath "$.requestId" exists
jsonpath "$.status" == "error"
jsonpath "$.message" == "Missing emailColumn parameter"

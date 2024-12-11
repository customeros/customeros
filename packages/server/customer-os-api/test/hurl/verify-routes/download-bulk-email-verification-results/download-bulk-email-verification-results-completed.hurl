# Test: Download Bulk Email Verification Results - completed
POST {{cos_url}}/verify/v1/email/bulk
Content-Type: multipart/form-data
X-CUSTOMER-OS-API-KEY: {{api_key}}
[MultipartFormData]
file: file,download-bulk-email-verification-results-completed.csv;

HTTP 200
[Captures]
job_id: jsonpath "$.jobId"
[Asserts]
jsonpath "$.message" == "File uploaded successfully"

GET {{cos_url}}/verify/v1/email/bulk/results/{{job_id}}/download
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[Options]
delay: 30000

HTTP 200
[Asserts]
header "Content-Type" contains "text/csv"
header "Content-Description" == "File Transfer"
header "Content-Transfer-Encoding" == "binary"
header "Content-Disposition" contains "attachment; filename="
body contains "Email,SyntaxIsValid,User,Domain,CleanEmail,IsFirewalled,Provider,SecureGatewayProvider,IsCatchAll,CanConnectSMTP,HasMXRecord,HasSPFRecord,TLSRequired,IsPrimaryDomain,PrimaryDomain,SkippedValidation,Deliverable,IsMailboxFull,IsRoleAccount,IsSystemGenerated,IsFreeAccount,SmtpSuccess,RetryValidation,TLSRequired,AlternateEmail"
body contains "silviu@openline.ai,true,silviu,openline.ai,silviu@openline.ai,false,google workspace"
body contains "silviu@customeros.ai,true,silviu,customeros.ai,silviu@customeros.ai,false,google workspace"
body contains "dyptorden@gmail.com,true,dyptorden,gmail.com,dyptorden@gmail.com,false,google workspace"
body contains "silviu,false,,,,false,,,false,false,false,false,false,false,,false,,false,false,false,false,false,false,false"
body contains "test2@example.com,true,test2,example.com,test2@example.com,false"
body contains "dyptorden@yahoo.com,true,dyptorden,yahoo.com,dyptorden@yahoo.com,false,yahoo"

# Test: Get Invoices For Existing Organization
GET {{cos_url}}/billing/v1/organizations/f75add8c-4f00-442d-8bc5-e739a655baff/invoices
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.invoices" count == 4
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')]" count == 1
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')].number" includes "PPG-04195"
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')].amount" includes 2.06
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')].currency" includes "USD"
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')].id" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')].dueDate" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'VOID')].publicUrl" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')]" count == 1
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')].number" includes "XGG-87291"
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')].amount" includes 2.06
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')].currency" includes "USD"
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')].id" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')].dueDate" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'PAID')].publicUrl" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')]" count == 1
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')].number" includes "ZXL-37829"
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')].amount" includes 2.06
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')].currency" includes "USD"
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')].id" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')].dueDate" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'OVERDUE')].publicUrl" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')]" count == 1
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')].number" includes "FMW-53533"
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')].amount" includes 2.06
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')].currency" includes "USD"
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')].id" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')].dueDate" exists
jsonpath "$.invoices[?(@.invoiceStatus == 'DUE')].publicUrl" exists

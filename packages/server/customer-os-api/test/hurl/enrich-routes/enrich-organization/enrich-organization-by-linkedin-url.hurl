# Test: Enrich Organization by linkedinUrl
GET {{cos_url}}/enrich/v1/organization
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
linkedinUrl: https://www.linkedin.com/company/customer-os/

HTTP *
[Asserts]
status == 200
header "Content-Type" == "application/json"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.data" exists
jsonpath "$.data.name" ==  "CustomerOS"
jsonpath "$.data.domain" == "customeros.ai"
jsonpath "$.data.website" == "https://customeros.ai"
jsonpath "$.data.description" exists
jsonpath "$.data.longDescription" exists
jsonpath "$.data.employees" exists
jsonpath "$.data.foundedYear"  == 2022
jsonpath "$.data.logos" isCollection
jsonpath "$.data.industry.industry" == "Technology, Information and Internet"
jsonpath "$.data.socials" isCollection
jsonpath "$.data.location.isHeadquarter" == true
jsonpath "$.data.location.country" == "United Kingdom"
jsonpath "$.data.location.countryCodeA2" == "GB"
jsonpath "$.data.location.city" == "London"

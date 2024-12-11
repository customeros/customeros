# Test: Successful IP intelligence lookup
GET {{cos_url}}/verify/v1/ip?address=8.8.8.8
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.status" == "success"
jsonpath "$.requestId" exists
jsonpath "$.ip.ipAddress" == "8.8.8.8"
# Verify threats object structure
jsonpath "$.ip.threats" isCollection
jsonpath "$.ip.threats.isProxy" == false
jsonpath "$.ip.threats.isVpn" == false
jsonpath "$.ip.threats.isTor" == false
jsonpath "$.ip.threats.isUnallocated" == false
jsonpath "$.ip.threats.isDatacenter" == false
jsonpath "$.ip.threats.isCloudRelay" == false
jsonpath "$.ip.threats.isMobile" == false
# Verify geolocation object structure
jsonpath "$.ip.geolocation" isCollection
jsonpath "$.ip.geolocation.city" == ""
jsonpath "$.ip.geolocation.country" == "United States"
jsonpath "$.ip.geolocation.countryIso" == "US"
jsonpath "$.ip.geolocation.isEuropeanUnion" == false
# Verify timezone object structure
jsonpath "$.ip.time_zone" isCollection
jsonpath "$.ip.time_zone.name" == ""
jsonpath "$.ip.time_zone.abbr" == ""
jsonpath "$.ip.time_zone.offset" == ""
jsonpath "$.ip.time_zone.is_dst" == false
jsonpath "$.ip.time_zone.current_time" isString
# Verify network object structure
jsonpath "$.ip.network" isCollection
jsonpath "$.ip.network.asn" == "AS15169"
jsonpath "$.ip.network.name" == "Google LLC"
jsonpath "$.ip.network.domain" == "google.com"
jsonpath "$.ip.network.type" == "business"
jsonpath "$.ip.organization" isCollection
jsonpath "$.ip.organization.name" == ""
jsonpath "$.ip.organization.domain" == ""
jsonpath "$.ip.organization.linkedin" == ""

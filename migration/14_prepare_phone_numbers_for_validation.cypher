# https://github.com/openline-ai/openline-customer-os/issues/857

#========== migrate e164 into rawPhoneNumber for non validated phone numbers

MATCH (p:PhoneNumber)
WHERE p.validated is null AND p.e164 is not null AND p.e164 <> ''
SET p.rawPhoneNumber = p.e164
REMOVE p.e164;
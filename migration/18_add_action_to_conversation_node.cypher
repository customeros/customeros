# replace <tenant> with the tenant name

MATCH (c:Conversation) SET c:Action;

# for each tenant, run the following query
MATCH (c:Conversation_<tenant>) SET c:Action_<tenant>;
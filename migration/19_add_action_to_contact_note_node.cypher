# replace <tenant> with the tenant name

MATCH (c:Contact)--(n:Note) SET n:Action;

# for each tenant, run the following query
MATCH (c:Contact_<tenant>)--(n:Note_<tenant>) SET n:Action_<tenant>;
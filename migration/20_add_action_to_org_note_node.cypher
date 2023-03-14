# replace <tenant> with the tenant name

MATCH (o:Organization)--(n:Note) SET n:Action;

# for each tenant, run the following query
MATCH (o:Organization_<tenant>)--(n:Note_<tenant>) SET n:Action_<tenant>;
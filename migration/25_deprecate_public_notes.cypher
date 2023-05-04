#Delete public notes
#remove public property

#verify before removing data
MATCH (n:Note_<TENANT_NAME_HERE>)--(:ExternalSystem {id:"zendesk_support"}) where n.public = true return count(n);
#remove data
MATCH (n:Note_<TENANT_NAME_HERE>)--(:ExternalSystem {id:"zendesk_support"}) where n.public = true detach delete n;

MATCH (n:Note_<TENANT_NAME_HERE>) where n.public = true remove n.public;
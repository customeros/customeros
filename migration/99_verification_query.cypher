# CHECK 1 - Result should be 0

CALL {
 MATCH (t1:Tenant)--(c:Contact)-[rel]-(tag:Tag)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(:Contact)-[rel]-(:Organization)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(:Contact)-[rel]-(:ExternalSystem)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(:Organization)-[rel]-(:ExternalSystem)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(:Note)-[rel]-(:ExternalSystem)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(c1:Contact)-[rel]-(:Conversation)--(c2:Contact)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(u1:User)-[rel]-(:Conversation)--(c2:Contact)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
 UNION
 MATCH (t1:Tenant)--(u1:User)-[rel]-(:Conversation)--(u2:User)--(t2:Tenant)
 WHERE t1.name <> t2.name
 return count(rel) as x
} return sum(x) as Problematic_relationships;

# CHECK 2 - Check amount of labels per node
CALL {
 MATCH (node:Tenant) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 1 return count(nodeCount) as x
 UNION
 MATCH (node:EmailDomain) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 1 return count(nodeCount) as x
 UNION
 MATCH (node:AlternateOrganization) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 1 return count(nodeCount) as x
 UNION
 MATCH (node:AlternateContact) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 1 return count(nodeCount) as x
 UNION
 MATCH (node:Tag) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:ExternalSystem) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:OrganizationType) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:User) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Contact) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Organization) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:JobRole) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Email) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Conversation) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Location) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Place) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:Note) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:PhoneNumber) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 2 return count(nodeCount) as x
 UNION
 MATCH (node:CustomField) with node, labels(node) as labs unwind labs as labsList with node, count(node) as nodeCount where nodeCount <> 3 return count(nodeCount) as x
} return sum(x) as Problematic_nodes;

# CHECK 3 - Parameterized query, to be used for each tenant

:param { tenant: "openline" };

CALL {
 MATCH (:Tenant {name:$tenant})--(n:Contact) WHERE NOT 'Contact_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:Organization) WHERE NOT 'Organization_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:Organization) WHERE NOT 'Organization_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:ExternalSystem) WHERE NOT 'ExternalSystem_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:OrganizationType) WHERE NOT 'OrganizationType_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:User) WHERE NOT 'User_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:Tag) WHERE NOT 'Tag_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})--(n:Email) WHERE NOT 'Email_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:Email) WHERE NOT 'Email_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:Conversation) WHERE NOT 'Conversation_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:JobRole) WHERE NOT 'JobRole_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:Location) WHERE NOT 'Location_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:CustomField) WHERE NOT 'CustomField_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:PhoneNumber) WHERE NOT 'PhoneNumber_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:Note) WHERE NOT 'Note_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(n:Action) WHERE NOT 'Action_'+$tenant  in labels(n) return count(n) as x
 UNION
 MATCH (:Tenant {name:$tenant})-[*2]-(:Location)--(n:Place) WHERE NOT 'Place_'+$tenant  in labels(n) return count(n) as x
} return sum(x) as Problematic_nodes;


# CHECK 4 - Query to verify labels mix
MATCH (n) RETURN count(labels(n)), labels(n);
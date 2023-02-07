# Below script to be executed on env with synced data when https://github.com/openline-ai/openline-customer-os/issues/692 is released

# Execute below script for each tenant
# Openline tenant
DROP INDEX basicSearchStandard_openline IF EXISTS;
CREATE FULLTEXT INDEX basicSearchStandard_openline FOR (n:Contact_openline|Email_openline|Organization_openline) ON EACH [n.firstName, n.lastName, n.name, n.email]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'standard',
    `fulltext.eventually_consistent`: true
  }
};

DROP INDEX basicSearchSimple_openline IF EXISTS;
CREATE FULLTEXT INDEX basicSearchSimple_openline FOR (n:Contact_openline|Email_openline|Organization_openline) ON EACH [n.firstName, n.lastName, n.email, n.name]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'simple',
    `fulltext.eventually_consistent`: true
  }
}

# Test tenant
DROP INDEX basicSearchStandard_test IF EXISTS;
CREATE FULLTEXT INDEX basicSearchStandard_test FOR (n:Contact_test|Email_test|Organization_test) ON EACH [n.firstName, n.lastName, n.name, n.email]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'standard',
    `fulltext.eventually_consistent`: true
  }
};

DROP INDEX basicSearchSimple_test IF EXISTS;
CREATE FULLTEXT INDEX basicSearchSimple_test FOR (n:Contact_test|Email_test|Organization_test) ON EACH [n.firstName, n.lastName, n.email, n.name]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'simple',
    `fulltext.eventually_consistent`: true
  }
}

# verify that the index is created
CALL {
    CALL db.index.fulltext.queryNodes("basicSearchStandard_openline", "matt~") YIELD node, score RETURN score, node, labels(node) as labels limit 50
    union
    CALL db.index.fulltext.queryNodes("basicSearchSimple_openline", "matt~") YIELD node, score RETURN score, node, labels(node) as labels limit 50
}
with labels, node, score order by score desc
with labels, node, collect(score) as scores
return labels, head(scores) as score, node order by score desc limit 50;
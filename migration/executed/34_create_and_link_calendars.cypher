// Antoine
MATCH (u:User {id: "da4aa95e-db0d-4fc4-8294-73b8c8a64847"})
MERGE (u)-[:HAS_CALENDAR]->(c:Calendar {id:randomUUID()})
  ON CREATE SET
  c.link = "https://cal.com/valot",
  c.calType = "CALCOM",
  c.primary = true,
  c.createdAt = datetime({timezone: "UTC"}),
  c.updatedAt = datetime({timezone: "UTC"}),
  c.source="openline",
  c.sourceOfTruth="openline",
  c.appSource = "manual",
  c:Calendar_agemoai
  ON MATCH SET c.link="https://cal.com/valot"

// Matt
MATCH (u:User {id: "3fdbc141-ffff-4f8c-98c9-1bcbf6b99bb8"})
MERGE (u)-[:HAS_CALENDAR]->(c:Calendar {id:randomUUID()})
  ON CREATE SET
  c.link = "https://cal.com/mbrown",
  c.calType = "CALCOM",
  c.primary = true,
  c.createdAt = datetime({timezone: "UTC"}),
  c.updatedAt = datetime({timezone: "UTC"}),
  c.source="openline",
  c.sourceOfTruth="openline",
  c.appSource = "manual",
  c:Calendar_agemoai
  ON MATCH SET c.link="https://cal.com/mbrown"

// Vasi
MATCH (u:User {id: "1638dd2b-a9f8-468b-8fc4-ed1d9f01b49d"})
MERGE (u)-[:HAS_CALENDAR]->(c:Calendar {id:randomUUID()})
  ON CREATE SET
  c.link = "https://cal.com/vasicoscotin",
  c.calType = "CALCOM",
  c.primary = true,
  c.createdAt = datetime({timezone: "UTC"}),
  c.updatedAt = datetime({timezone: "UTC"}),
  c.source="openline",
  c.sourceOfTruth="openline",
  c.appSource = "manual",
  c:Calendar_agemoai

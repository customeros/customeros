// Antoine
match (u:User_agemoai {id: "da4aa95e-db0d-4fc4-8294-73b8c8a64847"}) match (jr:JobRole_agemoai {id: "5937042d-142c-11ee-8b0c-fa8b76c30fd7"}) merge (u)-[:WORKS_AS]->(jr)
// Matt
match (u:User_agemoai {id: "3fdbc141-ffff-4f8c-98c9-1bcbf6b99bb8"}) match (jr:JobRole_agemoai {id: "3c73eafe-1420-11ee-8b0c-fa8b76c30fd7"}) merge (u)-[:WORKS_AS]->(jr)
// Vasi
match (u:User_agemoai {id: "1638dd2b-a9f8-468b-8fc4-ed1d9f01b49d"}) match (jr:JobRole_agemoai {id: "7e374585-142d-11ee-8b0c-fa8b76c30fd7"}) merge (u)-[:WORKS_AS]->(jr)
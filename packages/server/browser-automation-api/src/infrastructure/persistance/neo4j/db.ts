import neo4j from "neo4j-driver";

import { logger } from "@/infrastructure/logger";

export const db = neo4j.driver(
  process.env.NEO4J_URL ?? "",
  neo4j.auth.basic(process.env.NEO4J_USER ?? "", process.env.NEO4J_PASS ?? ""),
  { logging: { logger: (_, m) => logger.info(m) } },
);

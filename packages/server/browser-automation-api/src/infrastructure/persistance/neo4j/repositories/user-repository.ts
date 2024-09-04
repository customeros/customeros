import { db } from "../db";
import { ErrorParser } from "@/util/error";
import { logger } from "@/infrastructure/logger";

export class UserRepository {
  constructor() {}

  async getUserByEmail(tenant: string, email: string) {
    try {
      const session = db.session();
      const result = await session.run(
        `
            MATCH (:Tenant {name:$tenantName})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email)
       			WHERE e.email=$email OR e.rawEmail=$email
       			RETURN DISTINCT(u) ORDER by u.createdAt ASC limit 1
          `,
        {
          tenantName: tenant,
          email: email,
        },
      );

      const user = result.records?.[0]?.toObject()?.u?.properties;

      return user;
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in Neo4j UserRepository: ", {
        error: error.message,
        details: error.details,
      });
      throw error;
    }
  }
}

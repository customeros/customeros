import postgres from "postgres";
import { migrate } from "drizzle-orm/postgres-js/migrator";
import { drizzle } from "drizzle-orm/postgres-js";

import { connectionUrl } from "./config";

const pg = postgres(connectionUrl, { max: 1 });
const db = drizzle(pg);

(async () => {
  try {
    console.info(`Running migrations on ${connectionUrl}`);
    await migrate(db, {
      migrationsFolder: "./src/infrastructure/persistance/postgresql/drizzle",
    });
    await pg.end();
  } catch (err) {
    console.error(err);
  }
})();

import postgres from "postgres";
import { drizzle } from "drizzle-orm/postgres-js";

import { connectionUrl } from "./config";

const pg = postgres(connectionUrl, {
  ssl: {
    rejectUnauthorized: false,
  },
});
export const db = drizzle(pg);

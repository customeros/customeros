import postgres from "postgres";
import { drizzle } from "drizzle-orm/postgres-js";

import { connectionUrl } from "./config";

const pg = postgres(connectionUrl, {
  ssl: process.env.NODE_ENV !== "development",
});
export const db = drizzle(pg);

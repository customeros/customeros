import postgres from "postgres";
import { drizzle } from "drizzle-orm/postgres-js";

import { connectionUrl } from "./config";

const isDev = process.env.NODE_ENV !== "production";

const pg = postgres(connectionUrl, {
  ssl: !isDev
    ? {
        rejectUnauthorized: false,
      }
    : undefined,
});
export const db = drizzle(pg);

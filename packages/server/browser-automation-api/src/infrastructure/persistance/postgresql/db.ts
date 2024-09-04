import postgres from "postgres";
import { drizzle } from "drizzle-orm/postgres-js";

import { connectionUrl } from "./config";

const pg = postgres(connectionUrl);
export const db = drizzle(pg);

import { defineConfig } from "drizzle-kit";

const isDev = process.env.NODE_ENV === "development";

export default defineConfig({
  schema: "./src/infrastructure/persistance/postgresql/drizzle/schema.ts",
  out: "./src/infrastructure/persistance/postgresql/drizzle",
  dialect: "postgresql",
  dbCredentials: {
    host: process.env.POSTGRES_HOST ?? "",
    user: process.env.POSTGRES_USER ?? "",
    port: parseInt(process.env.POSTGRES_PORT ?? "5432"),
    password: process.env.POSTGRES_PASS ?? "",
    database: process.env.POSTGRES_NAME ?? "",
    ssl: !isDev
      ? {
          rejectUnauthorized: false,
        }
      : undefined,
  },
});

const USER = process.env.POSTGRES_USER ?? "";
const PASS = process.env.POSTGRES_PASS ?? "";
const HOST = process.env.POSTGRES_HOST ?? "";
const PORT = process.env.POSTGRES_PORT ?? "";
const DATABASE = process.env.POSTGRES_NAME ?? "";

export const connectionUrl = `postgresql://${USER}:${PASS}@${HOST}:${PORT}/${DATABASE}`;

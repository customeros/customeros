import winston from "winston";

const { combine, timestamp, printf, colorize } = winston.format;

const logFormat = printf(({ level, message, timestamp, details, error }) => {
  return `${timestamp} [${level}]: ${message} ${details ? "- " + details : ""} ${error ? "- " + error : ""}`;
});

export const logger = winston.createLogger({
  level: "info",
  format: combine(
    timestamp({ format: "YYYY-MM-DD HH:mm:ss" }),
    colorize(),
    logFormat,
  ),
  transports: [new winston.transports.Console()],
});

import winston from "winston";

const { combine, timestamp, printf, colorize } = winston.format;

const magenta = "\x1b[35m";
const reset = "\x1b[0m";

const logFormat = printf(
  ({ level, message, timestamp, details, error, source }) => {
    const coloredSource = source ? `${magenta}[${source}]${reset}` : "";
    return `${timestamp} [${level}]${coloredSource}: ${message} ${details ? "- " + details : ""} ${error ? "- " + error : ""}`;
  },
);

export const logger = winston.createLogger({
  level: "info",
  format: combine(
    timestamp({ format: "YYYY-MM-DD HH:mm:ss" }),
    colorize(),
    logFormat,
  ),
  transports: [new winston.transports.Console()],
});

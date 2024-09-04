import { logger } from "@/infrastructure/logger";

type ErrorSeverity = "low" | "medium" | "high" | "critical";
type ErrorCodes =
  | "UNKNOWN_ERROR"
  | "INTERNAL_ERROR"
  | "EXTERNAL_ERROR"
  | "VALIDATION_ERROR"
  | "INVARIANT_ERROR"
  | "APPLICATION_ERROR";

type ErrorDetails = {
  message: string;
  code: ErrorCodes;
  details?: string | null;
  severity: ErrorSeverity;
};

export class StandardError extends Error {
  code: ErrorCodes;
  timestamp: string;
  details: string | null;
  severity: ErrorSeverity;

  constructor({
    code,
    message,
    details = null,
    severity = "high",
  }: ErrorDetails) {
    super(message);
    this.code = code;
    this.details = details;
    this.severity = severity;
    this.timestamp = new Date().toISOString();
  }
}

export class ErrorParser {
  static parse(error: unknown) {
    if (error instanceof StandardError) {
      return error;
    }

    if (error instanceof Error) {
      return new StandardError({
        code: "INTERNAL_ERROR",
        message: error.message,
        details: error?.stack,
        severity: "critical",
      });
    }

    if ((error as any)?.response && (error as any)?.response?.status) {
      return new StandardError({
        code: "EXTERNAL_ERROR",
        message: (error as any)?.response?.statusText,
        details: JSON.stringify((error as any)?.response?.data),
        severity: "high",
      });
    }

    return new StandardError({
      code: "UNKNOWN_ERROR",
      message: "An unknown error occurred",
      details: JSON.stringify(error),
      severity: "high",
    });
  }
}

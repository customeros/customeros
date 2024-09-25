type ErrorSeverity = "low" | "medium" | "high" | "critical";
type ErrorType =
  | "UNKNOWN_ERROR"
  | "INTERNAL_ERROR"
  | "EXTERNAL_ERROR"
  | "EXTERNAL_ERROR"
  | "VALIDATION_ERROR"
  | "INVARIANT_ERROR"
  | "APPLICATION_ERROR";
type ErrorCode =
  // Session token is expired or invalid;
  | "S001"
  // Profile is already a connection;
  | "P001";

type ErrorDetails = {
  message: string;
  code: ErrorType;
  details?: string | null;
  severity: ErrorSeverity;
  reference?: ErrorCode | null;
};

export class StandardError extends Error {
  code: ErrorType;
  timestamp: string;
  details: string | null;
  severity: ErrorSeverity;
  reference: ErrorCode | null;

  constructor({
    code,
    message,
    details = null,
    severity = "high",
    reference = null,
  }: ErrorDetails) {
    super(message);
    this.code = code;
    this.details = details;
    this.severity = severity;
    this.reference = reference;
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

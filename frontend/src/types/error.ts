export type ErrorCode = 
  | "NOT_FOUND"
  | "VALIDATION"
  | "DATABASE"
  | "CONFLICT"
  | "FORBIDDEN"
  | "INTERNAL";

export interface AppError {
  code: ErrorCode;
  message: string;
  field?: string;
}

/**
 * Check if an error is an AppError
 */
export function isAppError(error: unknown): error is AppError {
  return (
    typeof error === "object" &&
    error !== null &&
    "code" in error &&
    "message" in error &&
    typeof (error as AppError).code === "string" &&
    typeof (error as AppError).message === "string"
  );
}

/**
 * Extract error message from various error types
 */
export function getErrorMessage(error: unknown): string {
  if (isAppError(error)) {
    return error.message;
  }
  
  if (error instanceof Error) {
    return error.message;
  }
  
  if (typeof error === "string") {
    return error;
  }
  
  return "An unexpected error occurred";
}

/**
 * Get user-friendly error message based on error code
 */
export function getFriendlyErrorMessage(error: AppError): string {
  switch (error.code) {
    case "NOT_FOUND":
      return error.message;
    case "VALIDATION":
      return error.field 
        ? `${error.field}: ${error.message}`
        : error.message;
    case "DATABASE":
      return "A database error occurred. Please try again.";
    case "CONFLICT":
      return error.message;
    case "FORBIDDEN":
      return error.message;
    case "INTERNAL":
      return "An internal error occurred. Please try again later.";
    default:
      return error.message;
  }
}

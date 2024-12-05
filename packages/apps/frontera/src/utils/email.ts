/**
 * Validates if a string is a valid email local-part (before the '@').
 * @param localPart - The local-part of an email to validate.
 * @returns `true` if the local-part is valid, otherwise `false`.
 */
export function validateEmailLocalPart(localPart: string): boolean {
  // Regular expression for valid email local-part characters
  const localPartRegex = /^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$/;

  // Validate length and regex pattern
  return (
    localPart.length > 0 && // Must not be empty
    localPart.length <= 64 && // Must not exceed 64 characters
    !localPart.startsWith('.') && // Must not start with a dot
    !localPart.endsWith('.') && // Must not end with a dot
    !localPart.includes('..') && // Must not have consecutive dots
    localPartRegex.test(localPart) // Matches allowed characters
  );
}

export const validateEmail = (value: string): undefined | string => {
  if (value) {
    if (!/@/.test(value)) {
      return 'This email address needs an @ sign';
    }

    const domainCheck = /\.[a-z]{2,255}$/i;

    if (!domainCheck.test(value)) {
      return 'This email address needs a domain';
    }

    if (value?.startsWith('.') || value?.endsWith('.')) {
      return 'The email address cannot start or end with a dot';
    }

    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!regex.test(value)) {
      if (value?.includes('..')) {
        return 'The email address cannot contain consecutive dots';
      }

      return 'The email address contains invalid characters';
    }
  }

  return undefined;
};

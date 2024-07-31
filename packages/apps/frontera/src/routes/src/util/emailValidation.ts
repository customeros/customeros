export const validateEmail = (value: string): undefined | string => {
  if (value) {
    if (!/@/.test(value)) {
      return 'This email address needs an @ sign';
    }

    const domainCheck = /\.[a-z]{2,3}$/i;

    if (!domainCheck.test(value)) {
      return 'This email address needs a domain';
    }

    if (value.startsWith('.') || value.endsWith('.')) {
      return 'The email address cannot start or end with a dot';
    }

    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!regex.test(value)) {
      if (value.includes('..')) {
        return 'The email address cannot contain consecutive dots';
      }

      return 'The email address contains invalid characters';
    }
  }

  return undefined;
};

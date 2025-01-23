/**
 * Validates if a string contains only characters allowed in URL domains.
 * @param domain - The domain string to validate.
 * @returns `true` if the domain is valid, otherwise `false`.
 */
export function validateDomain(domain: string): boolean {
  // Regular expression to match valid URL domain names
  const domainRegex =
    /^(?!-)(?!.*\.-)(?!.*-\.)[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$/;

  // Test if the domain matches the regex
  return domainRegex.test(domain);
}

/**
 * Validates if a string is a valid URL with or without protocol and subdomain.
 * @param url - The URL string to validate.
 * @returns `true` if the URL is valid, otherwise `false`.
 */
export function validateUrl(url: string): boolean {
  // Regular expression to match valid URLs
  const urlRegex =
    /^(?:(https?:\/\/)?([a-zA-Z0-9-]+\.)?[a-zA-Z0-9-]+\.[a-zA-Z]{2,})(\/.*)?$/;

  // Test if the URL matches the regex
  return urlRegex.test(url);
}

/**
 * Validates if a string is a valid URL that must include a protocol.
 * @param url - The URL string to validate.
 * @returns `true` if the URL is valid and includes a protocol, otherwise `false`.
 */
export function validateFullUrl(url: string): boolean {
  // Regular expression to match valid URLs with protocol, optional port, query parameters, and fragments
  const urlWithProtocolRegex =
    /^(https?:\/\/)[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(:\d+)?(\/.*)?(\?.*)?(#.*)?$/;

  // Test if the URL matches the regex
  return urlWithProtocolRegex.test(url);
}

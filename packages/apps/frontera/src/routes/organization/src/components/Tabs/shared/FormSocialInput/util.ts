export function isKnownUrl(input = '') {
  const url = input.trim().toLowerCase();
  if (url.includes('twitter')) return 'twitter';
  if (url.includes('linkedin')) return 'linkedin';
}

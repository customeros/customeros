export function isKnownUrl(input = '') {
  const url = input.trim().toLowerCase();

  if (url.includes('twitter')) return 'twitter';
  if (url.includes('linkedin')) return 'linkedin';
  if (url.includes('facebook')) return 'facebook';
  if (url.includes('github')) return 'github';
  if (url.includes('instagram')) return 'instagram';
}

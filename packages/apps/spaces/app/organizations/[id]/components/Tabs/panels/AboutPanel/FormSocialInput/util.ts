export function isKnownUrl(input: string) {
  const url = input.trim().toLowerCase();
  if (url.includes('twitter')) return 'twitter';
  if (url.includes('linkedin')) return 'linkedin';
}

export function formatSocialUrl(value: string) {
  let url = value;

  if (url.startsWith('http')) {
    url = url.replace('https://', '');
  }
  if (url.startsWith('www')) {
    url = url.replace('www.', '');
  }
  if (url.includes('twitter')) {
    url = url.replace('twitter.com', '');
  }
  if (url.includes('linkedin')) {
    url = url.replace('linkedin.com/in', '');
  }

  return url;
}

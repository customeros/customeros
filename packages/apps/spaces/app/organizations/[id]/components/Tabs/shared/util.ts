export function formatSocialUrl(value = '') {
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
  if (url.includes('linkedin.com/in')) {
    url = url.replace('linkedin.com/in', '');
  }
  if (url.includes('linkedin.com/company')) {
    url = url.replace('linkedin.com/company', '');
  }

  return url;
}

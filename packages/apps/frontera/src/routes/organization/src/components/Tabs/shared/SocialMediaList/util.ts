export function isKnownUrl(input = '') {
  const url = input.trim().toLowerCase();

  if (url.includes('google')) return 'google';
  if (url.includes('instagram')) return 'instagram';
  if (url.includes('reddit')) return 'reddit';
  if (url.includes('snapchat')) return 'snapchat';
  if (url.includes('twitter')) return 'twitter';
  if (url.includes('discord')) return 'discord';
  if (url.includes('linkedin')) return 'linkedin';
  if (url.includes('tiktok')) return 'tiktok';
  if (url.includes('telegram')) return 'telegram';
  if (url.includes('clubhouse')) return 'clubhouse';
  if (url.includes('youtube')) return 'youtube';
  if (url.includes('pinterest')) return 'pinterest';
  if (url.includes('angellist')) return 'angellist';
  if (url.includes('github')) return 'github';
  if (url.includes('facebook')) return 'facebook';
  if (url.includes('slack')) return 'slack';
  if (url.includes('notion')) return 'notion';
}

export function isKnownUrl(input = '') {
  const url = input.trim().toLowerCase();

  if (url.includes('google.com')) return 'google';
  if (url.includes('instagram.com')) return 'instagram';
  if (url.includes('reddit.com')) return 'reddit';
  if (url.includes('snapchat.com')) return 'snapchat';
  if (url.includes('twitter.com')) return 'twitter';
  if (url.includes('discord.com')) return 'discord';
  if (url.includes('linkedin.com')) return 'linkedin';
  if (url.includes('tiktok.com')) return 'tiktok';
  if (url.includes('telegram.com')) return 'telegram';
  if (url.includes('clubhouse.com')) return 'clubhouse';
  if (url.includes('youtube.com')) return 'youtube';
  if (url.includes('pinterest.com')) return 'pinterest';
  if (url.includes('angellist.com')) return 'angellist';
  if (url.includes('github.com')) return 'github';
  if (url.includes('facebook.com')) return 'facebook';
  if (url.includes('slack.com')) return 'slack';
  if (url.includes('notion.com')) return 'notion';
}

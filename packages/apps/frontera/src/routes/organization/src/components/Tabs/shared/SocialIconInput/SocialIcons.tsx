import { X } from '@ui/media/logos/X';
import { Slack } from '@ui/media/logos/Slack';
import { Reddit } from '@ui/media/logos/Reddit';
import { Tiktok } from '@ui/media/logos/Tiktok';
import { Google } from '@ui/media/logos/Google';
import { Discord } from '@ui/media/logos/Discord';
import { Youtube } from '@ui/media/logos/Youtube';
import { Notion } from '@ui/media/logos/Notion.tsx';
import { Github } from '@ui/media/logos/Github.tsx';
import { Linkedin } from '@ui/media/logos/Linkedin';
import { Snapchat } from '@ui/media/logos/Snapchat';
import { Telegram } from '@ui/media/logos/Telegram';
import { Clubhouse } from '@ui/media/logos/Clubhouse';
import { Pinterest } from '@ui/media/logos/Pinterest';
import { Angellist } from '@ui/media/logos/Angellist';
import { Facebook } from '@ui/media/logos/Facebook.tsx';
import { Instagram } from '@ui/media/logos/Instagram.tsx';

import { isKnownUrl } from './util';

export const SocialIcon = ({
  children,
  url,
}: React.PropsWithChildren<{ url: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter') return <X className='size-4' />;
  if (knownUrl === 'facebook') return <Facebook className='size-4' />;
  if (knownUrl === 'linkedin') return <Linkedin className='size-4' />;
  if (knownUrl === 'github') return <Github className='size-4' />;
  if (knownUrl === 'instagram') return <Instagram className='size-4' />;
  if (knownUrl === 'youtube') return <Youtube className='size-4' />;
  if (knownUrl === 'pinterest') return <Pinterest className='size-4' />;
  if (knownUrl === 'angellist') return <Angellist className='size-4' />;
  if (knownUrl === 'notion') return <Notion className='size-4' />;
  if (knownUrl === 'clubhouse') return <Clubhouse className='size-4' />;
  if (knownUrl === 'discord') return <Discord className='size-4' />;
  if (knownUrl === 'slack') return <Slack className='size-4' />;
  if (knownUrl === 'tiktok') return <Tiktok className='size-4' />;
  if (knownUrl === 'telegram') return <Telegram className='size-4' />;
  if (knownUrl === 'snapchat') return <Snapchat className='size-4' />;
  if (knownUrl === 'reddit') return <Reddit className='size-4' />;
  if (knownUrl === 'google') return <Google className='size-4' />;

  return <>{children}</>;
};

import { cn } from '@ui/utils/cn';
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
  className,
  url,
}: React.PropsWithChildren<{ url: string; className?: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter') return <X className={cn('size-4', className)} />;
  if (knownUrl === 'facebook')
    return <Facebook className={cn('size-4', className)} />;
  if (knownUrl === 'linkedin')
    return <Linkedin className={cn('size-4', className)} />;
  if (knownUrl === 'github')
    return <Github className={cn('size-4', className)} />;
  if (knownUrl === 'instagram')
    return <Instagram className={cn('size-4', className)} />;
  if (knownUrl === 'youtube')
    return <Youtube className={cn('size-4', className)} />;
  if (knownUrl === 'pinterest')
    return <Pinterest className={cn('size-4', className)} />;
  if (knownUrl === 'angellist')
    return <Angellist className={cn('size-4', className)} />;
  if (knownUrl === 'notion')
    return <Notion className={cn('size-4', className)} />;
  if (knownUrl === 'clubhouse')
    return <Clubhouse className={cn('size-4', className)} />;
  if (knownUrl === 'discord')
    return <Discord className={cn('size-4', className)} />;
  if (knownUrl === 'slack')
    return <Slack className={cn('size-4', className)} />;
  if (knownUrl === 'tiktok')
    return <Tiktok className={cn('size-4', className)} />;
  if (knownUrl === 'telegram')
    return <Telegram className={cn('size-4', className)} />;
  if (knownUrl === 'snapchat')
    return <Snapchat className={cn('size-4', className)} />;
  if (knownUrl === 'reddit')
    return <Reddit className={cn('size-4', className)} />;
  if (knownUrl === 'google')
    return <Google className={cn('size-4', className)} />;

  return <>{children}</>;
};

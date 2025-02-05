import { cn } from '@ui/utils/cn';
import { Logo } from '@ui/media/Logo/Logo';
import { Slack } from '@ui/media/logos/Slack';
import { Reddit } from '@ui/media/logos/Reddit';
import { Tiktok } from '@ui/media/logos/Tiktok';
import { Google } from '@ui/media/logos/Google';
import { Discord } from '@ui/media/logos/Discord';
import { Notion } from '@ui/media/logos/Notion.tsx';
import { Snapchat } from '@ui/media/logos/Snapchat';
import { Telegram } from '@ui/media/logos/Telegram';
import { Clubhouse } from '@ui/media/logos/Clubhouse';
import { Pinterest } from '@ui/media/logos/Pinterest';
import { Angellist } from '@ui/media/logos/Angellist';

import { isKnownUrl } from './util';

export const SocialIcon = ({
  children,
  className,
  url,
}: React.PropsWithChildren<{ url: string; className?: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter')
    return <Logo name='twitter' className={cn(className)} />;
  if (knownUrl === 'facebook')
    return <Logo name='facebook' className={cn(className)} />;
  if (knownUrl === 'linkedin')
    return <Logo fill='#0A66C2' name='linkedin' className={cn(className)} />;
  if (knownUrl === 'github')
    return <Logo name='github' className={cn(className)} />;
  if (knownUrl === 'instagram')
    return <Logo name='instagram' className={cn(className)} />;
  if (knownUrl === 'youtube')
    return <Logo name='youtube' className={cn(className)} />;
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

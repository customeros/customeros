import { cn } from '@ui/utils/cn';
import { Logo } from '@ui/media/Logo/Logo';

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
  if (knownUrl === 'google')
    return <Logo name='google' className={cn('size-4', className)} />;

  if (children) {
    return <>{children}</>;
  }

  return <Logo name='default' className={cn(className)} />;
};

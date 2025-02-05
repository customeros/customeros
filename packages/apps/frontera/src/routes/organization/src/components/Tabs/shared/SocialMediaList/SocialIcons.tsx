import { cn } from '@ui/utils/cn';
import { Logo, LogoName } from '@ui/media/Logo/Logo';

import { isKnownUrl } from './util';

const logoConfig: Partial<Record<string, { fill?: string; name: LogoName }>> = {
  twitter: { name: 'twitter' },
  facebook: { name: 'facebook' },
  linkedin: { name: 'linkedin', fill: '#0A66C2' },
  github: { name: 'github' },
  instagram: { name: 'instagram' },
  youtube: { name: 'youtube' },
  google: { name: 'google' },
};
interface SocialIconProps {
  url: string;
  className?: string;
  children?: React.ReactNode;
}

export const SocialIcon = ({ children, className, url }: SocialIconProps) => {
  const knownUrl = isKnownUrl(url);
  const logoProps = logoConfig[knownUrl as keyof typeof logoConfig];

  if (logoProps) {
    return <Logo {...logoProps} className={cn(className)} />;
  }

  if (children) {
    return <>{children}</>;
  }

  return <Logo name='default' className={cn(className)} />;
};

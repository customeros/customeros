import { isKnownUrl } from './util';
import { X } from '@ui/media/logos/X';
import { Linkedin } from '@ui/media/logos/Linkedin';

export const SocialIcon = ({
  children,
  url,
}: React.PropsWithChildren<{ url: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter') return <X />;
  if (knownUrl === 'linkedin') return <Linkedin />;
  return <>{children}</>;
};

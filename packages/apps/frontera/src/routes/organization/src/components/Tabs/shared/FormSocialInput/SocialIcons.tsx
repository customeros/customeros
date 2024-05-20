import { X } from '@ui/media/logos/X';
import { Github } from '@ui/media/logos/Github.tsx';
import { Linkedin } from '@ui/media/logos/Linkedin';
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

  return <>{children}</>;
};

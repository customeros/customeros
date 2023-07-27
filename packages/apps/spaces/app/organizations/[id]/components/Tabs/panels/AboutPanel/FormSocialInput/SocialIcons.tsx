import { Icons } from '@ui/media/Icon';
import { isKnownUrl } from './util';
import Image from 'next/image';
import React from 'react';

export const SocialIcon = ({
  children,
  url,
}: React.PropsWithChildren<{ url: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter')
    return (
      <Image
        src={'/logos/twitterX.webp'}
        alt='Twitter'
        width={32}
        height={32}
      />
    );
  if (knownUrl === 'linkedin')
    return <Icons.Linkedin viewBox='0 0 32 32' strokeWidth='0' />;
  return <>{children}</>;
};

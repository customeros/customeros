import { Icons } from '@ui/media/Icon';

import { isKnownUrl } from './util';

export const SocialIcon = ({
  children,
  url,
}: React.PropsWithChildren<{ url: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter')
    return <Icons.Twitter viewBox='0 0 32 32' strokeWidth='0' />;
  if (knownUrl === 'linkedin')
    return <Icons.Linkedin viewBox='0 0 32 32' strokeWidth='0' />;
  return <>{children}</>;
};

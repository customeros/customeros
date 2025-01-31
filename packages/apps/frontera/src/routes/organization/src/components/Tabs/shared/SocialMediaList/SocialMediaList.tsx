import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { Social } from '@shared/types/__generated__/graphql.types.ts';

import { SocialMediaItem } from './SocialMediaItem.tsx';

interface SocialMediaListProps {
  dataTest?: string;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
}

export const SocialMediaList = observer(
  ({ isReadOnly, dataTest, leftElement }: SocialMediaListProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.getById(store.ui.focusRow ?? id);

    if (!organization || !organization?.value) return null;

    const filteredSocials = filterUniqueSocials(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      organization.value.socialMedia as any,
    );
    const socialOptions = filteredSocials.map((social) => ({
      value: social.id,
      label: social.url,
    }));

    return (
      <div className='flex flex-wrap gap-2'>
        {socialOptions.map(({ value: v, label: l }: SelectOption) => (
          <div key={v} className='w-auto '>
            <SocialMediaItem
              id={v}
              value={l}
              dataTest={dataTest}
              isReadOnly={isReadOnly}
              leftElement={leftElement}
              organization={organization}
            />
          </div>
        ))}
      </div>
    );
  },
);

//this function will be moved in the BE in the future

const getSocialDomain = (url: string) => {
  if (!url || typeof url !== 'string') return null;

  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = `https://${url}`;
  }

  const domainMap = {
    twitter: 'twitter',
    facebook: 'facebook',
    linkedin: 'linkedin',
    github: 'github',
    instagram: 'instagram',
    youtube: 'youtube',
    pinterest: 'pinterest',
    angel: 'angellist',
    notion: 'notion',
    clubhouse: 'clubhouse',
    discord: 'discord',
    slack: 'slack',
    tiktok: 'tiktok',
    telegram: 'telegram',
    snapchat: 'snapchat',
    reddit: 'reddit',
    google: 'google',
  };

  const hostname = new URL(url).hostname;

  return Object.keys(domainMap).find((key) => hostname.includes(key)) || null;
};

const filterUniqueSocials = (socialArray = []) => {
  const uniqueSocials = new Map();

  socialArray.forEach((social: Social) => {
    if (!social || typeof social !== 'object' || !social.url) return;
    const domain = getSocialDomain(social.url);

    if (domain && !uniqueSocials.has(domain)) {
      uniqueSocials.set(domain, social);
    }
  });

  return Array.from(uniqueSocials.values());
};

import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { Social } from '@shared/types/__generated__/graphql.types.ts';

import { isKnownUrl } from './util.ts';
import { SocialMediaItem } from './SocialMediaItem.tsx';

interface SocialMediaListProps {
  dataTest?: string;
  isReadOnly?: boolean;
}

export const SocialMediaList = observer(
  ({ isReadOnly, dataTest }: SocialMediaListProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.getById(id ?? store.ui.focusRow);

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
      <div className='flex flex-wrap gap-3'>
        {socialOptions.map(({ value: v, label: l }: SelectOption) => (
          <div key={v} className='w-fit '>
            <SocialMediaItem
              id={v}
              value={l}
              dataTest={dataTest}
              isReadOnly={isReadOnly}
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

  return isKnownUrl(url.trim().toLowerCase()) || null;
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

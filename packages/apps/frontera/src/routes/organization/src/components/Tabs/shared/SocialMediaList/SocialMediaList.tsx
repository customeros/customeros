import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { SocialMediaItem } from './SocialMediaItem.tsx';

interface SocialMediaListProps {
  dataTest?: string;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
  value?: { label: string; value: string }[];
}

export const SocialMediaList = observer(
  ({ value, isReadOnly, dataTest, leftElement }: SocialMediaListProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.getById(store.ui.focusRow ?? id);

    if (!organization || !organization?.value) return null;

    return (
      <div className='flex flex-wrap gap-2'>
        {value?.map(({ value: v, label: l }) => (
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

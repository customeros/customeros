import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit01 } from '@ui/media/icons/Edit01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

interface OpportunityNameProps {
  opportunityId: string;
}

export const OpportunityName = observer(
  ({ opportunityId }: OpportunityNameProps) => {
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);
    const opportunityStore = store.opportunities.value.get(opportunityId);

    const opportunityName = opportunityStore?.value.name;

    return (
      <div
        className='flex items-center gap-1'
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        onDoubleClick={() => {
          store.ui.commandMenu.setType('RenameOpportunityName');
          store.ui.commandMenu.setOpen(true);
        }}
      >
        <p>{opportunityName}</p>
        {isHovered && (
          <IconButton
            size='xxs'
            className=' '
            variant='ghost'
            aria-label='Edit tags'
            icon={<Edit01 className='text-gray-500' />}
            onClick={() => {
              store.ui.commandMenu.setType('RenameOpportunityName');
              store.ui.commandMenu.setOpen(true);
            }}
          />
        )}
      </div>
    );
  },
);

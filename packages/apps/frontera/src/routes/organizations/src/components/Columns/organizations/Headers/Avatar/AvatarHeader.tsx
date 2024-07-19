import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { CreateNewOrganizationModal } from '@organizations/components/shared/CreateNewOrganizationModal.tsx';

export const AvatarHeader = observer(() => {
  const store = useStore();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  return (
    <div className='flex w-[24px] items-center justify-center'>
      <Tooltip
        label='Create an organization'
        side='bottom'
        align='center'
        className={cn(enableFeature ? 'visible' : 'hidden')}
        asChild
      >
        <IconButton
          className={cn('size-6', enableFeature ? 'visible' : 'hidden')}
          size='xxs'
          variant='ghost'
          aria-label='create organization'
          data-test='create-organization-from-table'
          onClick={() => {
            store.ui.setIsEditingTableCell(true);
            setIsCreateModalOpen(true);
          }}
          icon={<Plus className='text-gray-400 size-5' />}
        />
      </Tooltip>
      <CreateNewOrganizationModal
        isOpen={isCreateModalOpen}
        setIsOpen={setIsCreateModalOpen}
      />
    </div>
  );
});

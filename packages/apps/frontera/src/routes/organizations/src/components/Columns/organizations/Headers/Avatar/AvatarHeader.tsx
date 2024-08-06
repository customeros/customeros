import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';

export const AvatarHeader = observer(() => {
  const store = useStore();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');

  return (
    <div className='flex w-[24px] items-center justify-center'>
      <Tooltip
        asChild
        side='bottom'
        align='center'
        label='Create an organization'
        className={cn(enableFeature ? 'visible' : 'hidden')}
      >
        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='create organization'
          data-test='create-organization-from-table'
          icon={<Plus className='text-gray-400 size-5' />}
          className={cn('size-6', enableFeature ? 'visible' : 'hidden')}
          onClick={() => {
            store.ui.commandMenu.setType('AddNewOrganization');
            store.ui.commandMenu.setOpen(true);
          }}
        />
      </Tooltip>
    </div>
  );
});

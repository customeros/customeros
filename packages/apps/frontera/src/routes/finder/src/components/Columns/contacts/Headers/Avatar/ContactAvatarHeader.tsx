import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';

export const ContactAvatarHeader = observer(() => {
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const store = useStore();

  return (
    <div className='flex w-[26px] items-center justify-center'>
      <Tooltip
        asChild
        side='bottom'
        align='center'
        label='Create contact'
        className={cn(enableFeature ? 'visible' : 'hidden')}
      >
        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='create contact'
          data-test='create-contact-from-table'
          icon={<Plus className='text-gray-400 size-5' />}
          className={cn('size-6', enableFeature ? 'visible' : 'hidden')}
          onClick={() => {
            store.ui.commandMenu.setType('AddContactsBulk');
            store.ui.commandMenu.setOpen(true);
          }}
        />
      </Tooltip>
    </div>
  );
});

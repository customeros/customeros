import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';

export const ContactAvatarHeader = observer(() => {
  const enableFeature = useFeatureIsOn('gp-dedicated-1');

  return (
    <div className='flex w-[24px] items-center justify-center'>
      <Tooltip
        label='Create contact'
        side='bottom'
        align='center'
        className={cn(enableFeature ? 'visible' : 'hidden')}
        asChild
      >
        <IconButton
          className={cn('size-6', enableFeature ? 'visible' : 'hidden')}
          size='xxs'
          isDisabled
          variant='ghost'
          aria-label='create contact'
          data-test='create-contact-from-table'
          onClick={() => null}
          icon={<Plus className='text-gray-400 size-5' />}
        />
      </Tooltip>
    </div>
  );
});

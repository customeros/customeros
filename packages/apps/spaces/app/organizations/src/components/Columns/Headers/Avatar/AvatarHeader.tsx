import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { useOrganizationsPageMethods } from '@organizations/hooks/useOrganizationsPageMethods';

export const AvatarHeader = () => {
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const { createOrganization } = useOrganizationsPageMethods();

  const handleCreateOrganization = () => {
    createOrganization.mutate({ input: { name: 'Unnamed' } });
  };

  return (
    <div className='flex w-[42px] items-center justify-center'>
      <Tooltip
        label='Create an organization'
        side='bottom'
        align='center'
        className={cn(enableFeature ? 'visible' : 'hidden')}
        asChild={false}
      >
        <IconButton
          className={cn(enableFeature ? 'visible' : 'hidden')}
          size='sm'
          variant='ghost'
          aria-label='create organization'
          onClick={handleCreateOrganization}
          isLoading={createOrganization.isPending}
          icon={<Plus className='text-gray-400 size-5' />}
        />
      </Tooltip>
    </div>
  );
};

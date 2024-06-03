import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import {
  OrganizationStage,
  OrganizationRelationship,
} from '@shared/types/__generated__/graphql.types';

export const AvatarHeader = observer(() => {
  const store = useStore();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [searchParams] = useSearchParams();

  const preset = searchParams?.get('preset');

  const tableViewName = store.tableViewDefs.getById(`${preset}`)?.value.name;

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
          onClick={() =>
            store.organizations.create(
              defaultValuesNewOrganization(tableViewName ?? ''),
            )
          }
          icon={<Plus className='text-gray-400 size-5' />}
        />
      </Tooltip>
    </div>
  );
});

const defaultValuesNewOrganization = (organizationName: string) => {
  switch (organizationName) {
    case 'Customers':
      return {
        relationship: OrganizationRelationship.Customer,
        stage: OrganizationStage.Onboarding,
      };
    case 'Leads':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Lead,
      };
    case 'Nurture':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Target,
      };
    case 'All orgs':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Target,
      };
  }
};

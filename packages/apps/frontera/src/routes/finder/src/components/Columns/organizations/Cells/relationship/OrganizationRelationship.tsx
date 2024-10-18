import { useState } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { OrganizationStage, OrganizationRelationship } from '@graphql/types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { relationshipOptions } from './util';

interface OrganizationRelationshipProps {
  id: string;
  dataTest?: string;
}

export const OrganizationRelationshipCell = observer(
  ({ id, dataTest }: OrganizationRelationshipProps) => {
    const store = useStore();
    const [isEditing, setIsEditing] = useState(false);

    const organization = store.organizations.value.get(id);

    const enrichedOrg = organization?.value.enrichDetails;
    const enrichingStatus =
      !enrichedOrg?.enrichedAt &&
      enrichedOrg?.requestedAt &&
      !enrichedOrg?.failedAt;

    const value = relationshipOptions.find(
      (option) => option.value === organization?.value.relationship,
    );

    const handleSelect = (option: SelectOption<OrganizationRelationship>) => {
      if (!organization) return;

      organization.value.relationship = option.value;
      organization.value.stage = match(option.value)
        .with(OrganizationRelationship.Prospect, () => OrganizationStage.Lead)
        .with(
          OrganizationRelationship.Customer,
          () => OrganizationStage.InitialValue,
        )
        .with(
          OrganizationRelationship.NotAFit,
          () => OrganizationStage.Unqualified,
        )
        .with(
          OrganizationRelationship.FormerCustomer,
          () => OrganizationStage.Target,
        )
        .otherwise(() => undefined);

      organization.commit();
      setIsEditing(false);
    };

    return (
      <div className='flex gap-1 items-center group/relationship'>
        <p
          onDoubleClick={() => setIsEditing(true)}
          data-test='organization-relationship-in-all-orgs-table'
          className={cn(
            'cursor-default text-gray-700',
            !value && 'text-gray-400',
          )}
        >
          {value?.label
            ? value.label
            : enrichingStatus
            ? 'Enriching...'
            : 'Not set'}
        </p>
        <Menu open={isEditing} onOpenChange={setIsEditing}>
          <MenuButton asChild>
            <IconButton
              size='xxs'
              variant='ghost'
              id='edit-button'
              dataTest={dataTest}
              aria-label='edit relationship'
              onClick={() => setIsEditing(true)}
              icon={<Edit03 className='text-gray-500' />}
              className={cn(
                'rounded-md opacity-0 group-hover/relationship:opacity-100 min-w-5',
                isEditing && 'opacity-100',
              )}
            />
          </MenuButton>
          <MenuList onKeyDown={(e) => e.stopPropagation()}>
            {relationshipOptions
              .filter(
                (option) =>
                  !(
                    value?.label === 'Customer' && option.label === 'Prospect'
                  ) &&
                  !(
                    value?.label === 'Not a Fit' && option.label === 'Prospect'
                  ),
              )
              .map((option) => (
                <MenuItem
                  key={option.value.toString()}
                  onClick={() => handleSelect(option)}
                  data-test={`relationship-${option.value}`}
                >
                  {option.label}
                </MenuItem>
              ))}
          </MenuList>
        </Menu>
      </div>
    );
  },
);

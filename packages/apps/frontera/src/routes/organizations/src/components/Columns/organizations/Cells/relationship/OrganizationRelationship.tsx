import { useState } from 'react';

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

    const value = relationshipOptions.find(
      (option) => option.value === organization?.value.relationship,
    );

    const handleSelect = (option: SelectOption<OrganizationRelationship>) => {
      organization?.update((org) => {
        org.relationship = option.value;

        if (option.value === OrganizationRelationship.Prospect) {
          org.stage = OrganizationStage.Lead;
        }

        if (option.value === OrganizationRelationship.Customer) {
          org.stage = undefined;
        }

        if (option.value === OrganizationRelationship.NotAFit) {
          org.stage = OrganizationStage.Unqualified;
        }

        if (option.value === OrganizationRelationship.FormerCustomer) {
          org.stage = undefined;
        }

        return org;
      });
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
          {value?.label ?? 'No relationship'}
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
          <MenuList>
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

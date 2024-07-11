import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { OrganizationRelationship } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

import { relationshipOptions } from './util';

interface OrganizationRelationshipProps {
  id: string;
}

export const OrganizationRelationshipCell = observer(
  ({ id }: OrganizationRelationshipProps) => {
    const store = useStore();
    const [isEditing, setIsEditing] = useState(false);

    const organization = store.organizations.value.get(id);

    const value = relationshipOptions.find(
      (option) => option.value === organization?.value.relationship,
    );

    const handleSelect = (option: SelectOption<OrganizationRelationship>) => {
      organization?.update((org) => {
        org.relationship = option.value;

        return org;
      });
      setIsEditing(false);
    };

    return (
      <div className='flex gap-1 items-center group/relationship'>
        <p
          className={cn(
            'cursor-default text-gray-700',
            !value && 'text-gray-400',
          )}
          data-test='organization-relationship-in-all-orgs-table'
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.label ?? 'No relationship'}
        </p>
        <Menu open={isEditing} onOpenChange={setIsEditing}>
          <MenuButton asChild>
            <IconButton
              className={cn(
                'rounded-md opacity-0 group-hover/relationship:opacity-100',
                isEditing && 'opacity-100',
              )}
              aria-label='edit relationship'
              size='xxs'
              variant='ghost'
              id='edit-button'
              onClick={() => setIsEditing(true)}
              icon={<Edit03 className='text-gray-500' />}
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
                  data-test={`relationship-${option.label
                    .toLowerCase()
                    .replace(/\s+/g, '-')}`}
                  onClick={() => handleSelect(option)}
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

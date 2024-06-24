import { useParams } from 'react-router-dom';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding';
import { OrganizationRelationship } from '@graphql/types';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { relationshipOptions } from '@organizations/components/Columns/Cells/relationship/util';

const iconMap = {
  Customer: <ActivityHeart />,
  Prospect: <Seeding />,
  'Not a Fit': <MessageXCircle />,
  'Former Customer': <BrokenHeart />,
};

export const RelationshipButton = () => {
  const id = useParams()?.id as string;
  const store = useStore();
  const organization = store.organizations.value.get(id);

  const selectedValue = relationshipOptions.find(
    (option) => option.value === organization?.value?.relationship,
  );

  const spinnerColors =
    selectedValue?.value === OrganizationRelationship.Customer
      ? 'text-success-500 fill-succes-700'
      : 'text-gray-400 fill-gray-700';

  const iconTag = iconMap[selectedValue?.label as keyof typeof iconMap];

  return (
    <div>
      <Menu>
        <MenuButton asChild>
          <Tag
            variant='outline'
            size={'sm'}
            colorScheme={
              selectedValue?.value === OrganizationRelationship.Customer
                ? 'success'
                : 'gray'
            }
            className={cn(
              selectedValue?.value === OrganizationRelationship.Customer
                ? 'text-success-500'
                : 'text-gray-500',
              'rounded-full py-0.5 cursor-pointer',
            )}
          >
            <TagLeftIcon>
              {store.organizations.isLoading ? (
                <Spinner
                  label='Organization loading'
                  size='sm'
                  className={cn(spinnerColors)}
                />
              ) : (
                iconTag
              )}
            </TagLeftIcon>
            <TagLabel>{selectedValue?.label ?? 'Relationship'}</TagLabel>
          </Tag>
        </MenuButton>
        <MenuList className='min-w-[280px]'>
          {relationshipOptions
            .filter(
              (option) =>
                !(
                  selectedValue?.label === 'Customer' &&
                  option.label === 'Prospect'
                ) &&
                !(
                  selectedValue?.label === 'Not a Fit' &&
                  option.label === 'Prospect'
                ),
            )
            .map((option) => (
              <MenuItem
                key={option.value}
                onClick={() => {
                  organization?.update((prev) => ({
                    ...prev,
                    relationship: option.value,
                  }));
                }}
              >
                {iconMap[option.label as keyof typeof iconMap]}
                {option.label}
              </MenuItem>
            ))}
        </MenuList>
      </Menu>
    </div>
  );
};

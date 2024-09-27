import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { relationshipOptions } from '@finder/components/Columns/organizations/Cells/relationship/util';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { OrganizationStage, OrganizationRelationship } from '@graphql/types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

const iconMap = {
  Customer: <ActivityHeart />,
  Prospect: <Seeding />,
  'Not a Fit': <MessageXCircle />,
  'Former Customer': <BrokenHeart />,
};

export const RelationshipButton = observer(() => {
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
  };

  return (
    <div>
      <Menu>
        <MenuButton asChild data-test={`organization-account-relationship`}>
          <Tag
            size={'sm'}
            variant='outline'
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
                  size='sm'
                  label='Organization loading'
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
                onClick={() => handleSelect(option)}
                data-test={`relationship-${option.value}`}
              >
                {iconMap[option.label as keyof typeof iconMap]}
                {option.label}
              </MenuItem>
            ))}
        </MenuList>
      </Menu>
    </div>
  );
});

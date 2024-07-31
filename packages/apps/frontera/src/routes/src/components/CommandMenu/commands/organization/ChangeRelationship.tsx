import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding.tsx';
import { OrganizationRelationship } from '@graphql/types';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart.tsx';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart.tsx';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle.tsx';
import { relationshipOptions } from '@organization/components/Tabs/panels/AboutPanel/util.ts';
import {
  Command,
  CommandItem,
  CommandInput,
  useCommandState,
} from '@ui/overlay/CommandMenu';
const iconMap = {
  Customer: <ActivityHeart className='text-gray-500' />,
  Prospect: <Seeding className='text-gray-500' />,
  'Not a fit': <MessageXCircle className='text-gray-500' />,
  'Former Customer': <BrokenHeart className='text-gray-500' />,
};

export const ChangeRelationship = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = store.organizations.value.get(context.id as string);
  const label = `Organization - ${entity?.value?.name}`;

  const handleSelect = (value: OrganizationRelationship) => () => {
    if (!context.id) return;

    if (!entity) return;
    entity?.update((org) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      org.relationship = value;

      return org;
    });
    store.ui.commandMenu.toggle('ChangeRelationship');
  };

  const selectedRelationshipOption = relationshipOptions.find(
    (option) => option.value === entity?.value.relationship,
  );
  return (
    <Command label='Change Relationship'>
      <CommandInput label={label} placeholder='Change relationship...' />

      <Command.List>
        {relationshipOptions
          .filter(
            (option) =>
              !(
                selectedRelationshipOption?.label === 'Customer' &&
                option.label === 'Prospect'
              ) &&
              !(
                selectedRelationshipOption?.label === 'Not a fit' &&
                option.label === 'Prospect'
              ),
          )
          .map((option) => (
            <div
              className={
                cn(
                  (selectedRelationshipOption?.label === 'Customer' ||
                    selectedRelationshipOption?.label === 'Not a fit') &&
                    option.label === 'Prospect',
                ) && 'opacity-5 pointer-events-none'
              }
            >
              <CommandItem
                key={option.value}
                onSelect={handleSelect(option.value)}
                rightAccessory={
                  selectedRelationshipOption?.value === option.value ? (
                    <Check />
                  ) : null
                }
              >
                {iconMap[option.label as keyof typeof iconMap]}
                {option.label}
              </CommandItem>
            </div>
          ))}
      </Command.List>
    </Command>
  );
});

const EmptySearch = () => {
  const search = useCommandState((state) => state.search);

  return (
    <Command.Empty>{`No relationship status found with name "${search}".`}</Command.Empty>
  );
};

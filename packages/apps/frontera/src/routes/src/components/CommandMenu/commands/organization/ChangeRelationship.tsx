import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding.tsx';
import { OrganizationRelationship } from '@graphql/types';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart.tsx';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart.tsx';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { relationshipOptions } from '@organization/components/Tabs/panels/AboutPanel/util.ts';
const iconMap = {
  Customer: <ActivityHeart className='text-gray-500' />,
  Prospect: <Seeding className='text-gray-500' />,
  'Not a fit': <MessageXCircle className='text-gray-500' />,
  'Former Customer': <BrokenHeart className='text-gray-500' />,
};

export const ChangeRelationship = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = store.organizations.value.get((context.ids as string[])?.[0]);
  const label = `Organization - ${entity?.value?.name}`;

  const handleSelect = (value: OrganizationRelationship) => () => {
    if (!context.ids?.[0]) return;

    if (!entity) return;
    entity?.update((org) => {
      org.relationship = value;

      return org;
    });
    store.ui.commandMenu.toggle('ChangeRelationship');
  };

  const selectedRelationshipOption = relationshipOptions.find(
    (option) => option.value === entity?.value.relationship,
  );

  const options = relationshipOptions.filter(
    (option) =>
      !(
        selectedRelationshipOption?.label === 'Customer' &&
        option.label === 'Prospect'
      ) &&
      !(
        selectedRelationshipOption?.label === 'Not a fit' &&
        option.label === 'Prospect'
      ),
  );

  return (
    <Command label='Change Relationship'>
      <CommandInput label={label} placeholder='Change relationship...' />

      <Command.List>
        {options.map((option) => (
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
        ))}
      </Command.List>
    </Command>
  );
});

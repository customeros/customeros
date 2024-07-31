import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { User01 } from '@ui/media/icons/User01';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Calculator } from '@ui/media/icons/Calculator';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const OpportunityCommands = observer(() => {
  const store = useStore();
  const opportunity = store.opportunities.value.get(
    store.ui.commandMenu.context.id as string,
  );
  const label = `Opportunity - ${opportunity?.value.name}`;

  return (
    <Command>
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <CommandItem onSelect={() => {}} leftAccessory={<Columns03 />}>
          Change stage...
        </CommandItem>

        <CommandItem onSelect={() => {}} leftAccessory={<Calculator />}>
          Change ARR estimate
        </CommandItem>
        <CommandItem
          onSelect={() => {}}
          leftAccessory={<CurrencyDollarCircle />}
        >
          Change ARR currency...
        </CommandItem>
        <CommandItem onSelect={() => {}} leftAccessory={<Edit03 />}>
          Rename opportunity
        </CommandItem>
        <CommandItem
          leftAccessory={<User01 />}
          onSelect={() => {
            store.ui.commandMenu.setType('AssignOwner');
          }}
        >
          Assign owner...
        </CommandItem>
        <CommandItem onSelect={() => {}} leftAccessory={<Archive />}>
          Archive opportunity
        </CommandItem>
      </Command.List>
    </Command>
  );
});

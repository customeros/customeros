import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ArrowNarrowRight } from '@ui/media/icons/ArrowNarrowRight';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const GlobalHub = () => {
  return (
    <Command>
      <CommandInput placeholder='Type a command or search' />

      <Command.List>
        <GlobalSharedCommands />
      </Command.List>
    </Command>
  );
};

export const GlobalSharedCommands = observer(() => {
  const store = useStore();
  const navigate = useNavigate();

  const leadsPreset = store.tableViewDefs.leadsPreset;
  const targetsPreset = store.tableViewDefs.targetsPreset;
  const churnedPreset = store.tableViewDefs.churnedPreset;
  const customersPreset = store.tableViewDefs.defaultPreset;
  const addressBookPreset = store.tableViewDefs.addressBookPreset;

  const handleGoTo = (path: string, preset?: string) => {
    navigate(path + (preset ? `?preset=${preset}` : ''));
    store.ui.commandMenu.setOpen(false);
  };

  return (
    <>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/finder', leadsPreset)}
      >
        Go to Leads
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/finder', targetsPreset)}
      >
        Go to Targets
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/prospects')}
      >
        Go to Opportunities
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/finder', customersPreset)}
      >
        Go to Customers
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/finder', churnedPreset)}
      >
        Go to Former customers
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/finder', addressBookPreset)}
      >
        Go to Address book
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/settings')}
      >
        Go to Settings
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/customer-map')}
      >
        Go to Customer map
      </CommandItem>
    </>
  );
});

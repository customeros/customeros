import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ArrowNarrowRight } from '@ui/media/icons/ArrowNarrowRight';
import {
  Kbd,
  Command,
  CommandItem,
  CommandInput,
} from '@ui/overlay/CommandMenu';

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

  const targetsPreset = store.tableViewDefs.targetsPreset;
  const customersPreset = store.tableViewDefs.defaultPreset;
  const organizationsPreset = store.tableViewDefs.organizationsPreset;
  const contactsPreset = store.tableViewDefs.contactsPreset;
  const upcomingInvoicesPreset = store.tableViewDefs.upcomingInvoicesPreset;
  const contractsPreset = store.tableViewDefs.contractsPreset;

  const handleGoTo = (path: string, preset?: string) => {
    navigate(path + (preset ? `?preset=${preset}` : ''));
    store.ui.commandMenu.setOpen(false);
  };

  useEffect(() => {
    document.addEventListener('keydown', (e: KeyboardEvent) => {
      if (e.metaKey && e.key === 'k' && e.shiftKey) {
        store.ui.commandMenu.setType('GlobalHub');
        store.ui.commandMenu.setOpen(true);
      }
    });

    return () => {
      document.removeEventListener('keydown', (e: KeyboardEvent) => {
        if (e.metaKey && e.key === 'k' && e.shiftKey) {
          store.ui.commandMenu.setType('GlobalHub');
          store.ui.commandMenu.setOpen(true);
        }
      });
    };
  }, []);

  return (
    <>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_targets}
        rightAccessory={<KeyboardShortcut shortcut='T' />}
        onSelect={() => handleGoTo('/finder', targetsPreset)}
      >
        Go to Targets
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/prospects')}
        keywords={navigationKeywords.go_to_customers}
        rightAccessory={<KeyboardShortcut shortcut='O' />}
      >
        Go to Opportunities
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_customers}
        rightAccessory={<KeyboardShortcut shortcut='C' />}
        onSelect={() => handleGoTo('/finder', customersPreset)}
      >
        Go to Customers
      </CommandItem>

      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_address_book}
        rightAccessory={<KeyboardShortcut shortcut='Z' />}
        onSelect={() => handleGoTo('/finder', organizationsPreset)}
      >
        Go to Organizations
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_address_book}
        rightAccessory={<KeyboardShortcut shortcut='N' />}
        onSelect={() => handleGoTo('/finder', contactsPreset)}
      >
        Go to Contacts
      </CommandItem>

      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        rightAccessory={<KeyboardShortcut shortcut='I' />}
        keywords={navigationKeywords.go_to_scheduled_invoices}
        onSelect={() => handleGoTo('/finder', upcomingInvoicesPreset)}
      >
        Go to Invoices
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_contracts}
        rightAccessory={<KeyboardShortcut shortcut='R' />}
        onSelect={() => handleGoTo('/finder', contractsPreset)}
      >
        Go to Contracts
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/settings')}
        keywords={navigationKeywords.go_to_opportunities}
        rightAccessory={<KeyboardShortcut shortcut='S' />}
      >
        Go to Settings
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        onSelect={() => handleGoTo('/customer-map')}
        keywords={navigationKeywords.go_to_customer_map}
        rightAccessory={<KeyboardShortcut shortcut='D' />}
      >
        Go to Customer map
      </CommandItem>
    </>
  );
});

const KeyboardShortcut = ({ shortcut }: { shortcut: string }) => {
  return (
    <>
      <Kbd className='px-1.5'>G</Kbd>
      <span className='text-gray-500 text-[12px]'>then</span>
      <Kbd className='px-1.5'>{shortcut}</Kbd>
    </>
  );
};

const navigationKeywords = {
  go_to_contracts: ['go to', 'contracts', 'navigate'],
  go_to_contacts: ['go to', 'contacts', 'navigate', 'people'],
  go_to_targets: ['go to', 'targets', 'navigate'],
  go_to_customers: ['go to', 'customers', 'navigate'],
  go_to_address_book: ['go to', 'organizations', 'navigate'],
  go_to_opportunities: [
    'go to',
    'opportunities',
    'navigate',
    'deals',
    'pipeline',
  ],
  go_to_scheduled_invoices: [
    'go to',
    'invoices',
    'navigate',
    'past',
    'scheduled',
  ],
  go_to_settings: [
    'go to',
    'settings',
    'navigate',
    'accounts',
    'integrations',
    'apps',
    'emails',
    'billing',
    'data',
  ],
  go_to_customer_map: ['go to', 'customer', 'map', 'navigate', 'dashboard'],
};

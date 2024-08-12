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

  const leadsPreset = store.tableViewDefs.leadsPreset;
  const targetsPreset = store.tableViewDefs.targetsPreset;
  const churnedPreset = store.tableViewDefs.churnedPreset;
  const customersPreset = store.tableViewDefs.defaultPreset;
  const addressBookPreset = store.tableViewDefs.addressBookPreset;

  const handleGoTo = (path: string, preset?: string) => {
    navigate(path + (preset ? `?preset=${preset}` : ''));
    store.ui.commandMenu.setOpen(false);
  };

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.metaKey && e.key === 'k' && e.shiftKey) {
        store.ui.commandMenu.setType('GlobalHub');
        store.ui.commandMenu.setOpen(true);
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, []);

  return (
    <>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_leads}
        rightAccessory={<KeyboardShortcut shortcut='L' />}
        onSelect={() => handleGoTo('/finder', leadsPreset)}
      >
        Go to Leads
      </CommandItem>
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
        rightAccessory={<KeyboardShortcut shortcut='F' />}
        keywords={navigationKeywords.go_to_former_customers}
        onSelect={() => handleGoTo('/finder', churnedPreset)}
      >
        Go to Former customers
      </CommandItem>
      <CommandItem
        leftAccessory={<ArrowNarrowRight />}
        keywords={navigationKeywords.go_to_address_book}
        rightAccessory={<KeyboardShortcut shortcut='A' />}
        onSelect={() => handleGoTo('/finder', addressBookPreset)}
      >
        Go to Address book
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
  go_to_leads: ['go to', 'navigate', 'leads', 'prospect'],
  go_to_targets: ['go to', 'navigate', 'targets', 'prospect'],
  go_to_customers: ['go to', 'navigate', 'customers', 'relationship'],
  go_to_former_customers: [
    'go to',
    'navigate',
    'churned',
    'former customers',
    'relationship',
  ],
  go_to_address_book: [
    'go to',
    'navigate',
    'address book',
    'all contact',
    'all orgs',
    'leads',
    'targets',
    'customers',
    'former customers',
    'unqualified',
    'prospects',
  ],
  go_to_opportunities: [
    'go to',
    'navigate',
    'opportunities',
    'deals',
    'pipeline',
  ],
  go_to_my_portfolio: ['go to', 'navigate', 'my portfolio'],
  go_to_scheduled_invoices: [
    'go to',
    'navigate',
    'scheduled invoices',
    'past invoices',
    'billing',
  ],
  go_to_settings: [
    'go to',
    'navigate',
    'settings',
    'accounts',
    'integrations',
    'apps',
    'emails',
    'billing',
    'data',
  ],
  go_to_customer_map: [
    'go to',
    'navigate',
    'customer',
    'map',
    'dashboard',
    'charts',
    'graphs',
  ],
};

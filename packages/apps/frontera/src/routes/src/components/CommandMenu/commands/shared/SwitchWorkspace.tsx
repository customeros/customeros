import { useMemo, useState } from 'react';

import { useKey } from 'rooks';
import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { SwitchWorkspaceUsecase } from '@domain/usecases/settings-impersonate/switch-workspace.usecase.ts';

import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const SwitchWorkspace = observer(() => {
  const [searchTerm, setSearchTerm] = useState('');
  const store = useStore();

  const switchWorkspaceUsecase = useMemo(() => {
    return new SwitchWorkspaceUsecase();
  }, []);

  useKey('Escape', () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('GlobalHub');
  });

  return (
    <Command label='Switch workspace'>
      <CommandInput
        label={''}
        value={searchTerm}
        onValueChange={setSearchTerm}
        placeholder='Search a workspace'
        onKeyDownCapture={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            store.ui.commandMenu.setOpen(false);
          }
        }}
      />
      <CommandGroup>
        <Command.List>
          {store.common?.impersonateAccounts?.map((option) => (
            <CommandItem
              key={option?.tenant}
              onSelect={() => {
                switchWorkspaceUsecase.execute(option?.tenant);
              }}
              rightAccessory={
                option.tenant === store.session.value.tenant ? (
                  <Icon name={'check'} />
                ) : undefined
              }
            >
              {option?.tenant}
            </CommandItem>
          ))}
        </Command.List>
      </CommandGroup>
    </Command>
  );
});

import React, { ReactNode } from 'react';

import { Command, CommandInput } from '@ui/overlay/CommandMenu';
import { GlobalSharedCommands } from '@shared/components/CommandMenu/commands';

import { GlobalSearchResultNavigationCommands } from './GlobalSearchResultNavigationCommands.tsx';

export const CommandsContainer = ({
  children,
  label,
}: {
  label: string;
  children: ReactNode;
}) => {
  return (
    <Command>
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <Command.Group>{children}</Command.Group>
        <Command.Group heading='Navigate'>
          <GlobalSharedCommands />
        </Command.Group>

        <GlobalSearchResultNavigationCommands />
      </Command.List>
    </Command>
  );
};

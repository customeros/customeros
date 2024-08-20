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
    <Command
      filter={(value, search, keywords) => {
        const extendValue = value.replace(/\s/g, '') + keywords;
        const searchWithoutSpaces = search.replace(/\s/g, '');

        if (
          extendValue.toLowerCase().includes(searchWithoutSpaces.toLowerCase())
        )
          return 1;

        return 0;
      }}
    >
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <Command.Group>{children}</Command.Group>
        <GlobalSearchResultNavigationCommands />
        <Command.Group heading='Navigate'>
          <GlobalSharedCommands />
        </Command.Group>
      </Command.List>
    </Command>
  );
};

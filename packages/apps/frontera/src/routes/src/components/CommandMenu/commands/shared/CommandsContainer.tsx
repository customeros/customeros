import { ReactNode } from 'react';

import { Command, CommandInput } from '@ui/overlay/CommandMenu';
import { GlobalSharedCommands } from '@shared/components/CommandMenu/commands';

import { GlobalSearchResultNavigationCommands } from './GlobalSearchResultNavigationCommands.tsx';

export const CommandsContainer = ({
  children,
  label,
  dataTest,
}: {
  label: string;
  dataTest?: string;
  children: ReactNode;
}) => {
  return (
    <Command
      loop
      data-test={dataTest}
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
      <CommandInput
        label={label}
        data-test={`${dataTest}-input`}
        placeholder='Type a command or search'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <Command.Group>{children}</Command.Group>
        <GlobalSearchResultNavigationCommands />
        <Command.Group heading='Navigate' data-test={`${dataTest}-navigate`}>
          <GlobalSharedCommands dataTest={dataTest} />
        </Command.Group>
      </Command.List>
    </Command>
  );
};

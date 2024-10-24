import { Command } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

import { GlobalSharedCommands } from './GlobalHub';

export const ContactHub = () => {
  const label = `Contact`;

  return (
    <CommandsContainer label={label}>
      <Command.Group heading='Navigate'>
        <GlobalSharedCommands />
      </Command.Group>
    </CommandsContainer>
  );
};

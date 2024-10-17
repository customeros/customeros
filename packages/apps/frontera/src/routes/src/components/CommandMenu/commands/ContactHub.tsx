import { observer } from 'mobx-react-lite';

import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const ContactHub = observer(() => {
  const label = `Contact`;

  return (
    <CommandsContainer label={label}>
      <></>
    </CommandsContainer>
  );
});

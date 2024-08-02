import { useState } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText';

export const SetOpportunityNextSteps = observer(() => {
  const store = useStore();
  const [value, setValue] = useState('');
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(context.id as string);

  const label = match(context.entity)
    .with('Opportunity', () => `Opportunity - ${opportunity?.value?.name}`)
    .otherwise(() => 'Change ARR estimate');

  const handleEnterKey = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && e.metaKey) {
      opportunity?.update((o) => {
        const plainTextValue = value;

        o.nextSteps = plainTextValue;

        return o;
      });

      store.ui.commandMenu.setType('OpportunityCommands');
      store.ui.commandMenu.setOpen(false);
    }
  };

  return (
    <Command shouldFilter={false} onKeyDown={handleEnterKey}>
      <CommandInput asChild label={label} placeholder='Set next steps'>
        <Editor
          size='md'
          usePlainText
          namespace='opportunity-next-step'
          onChange={(html) => setValue(extractPlainText(html))}
        />
      </CommandInput>

      <Command.List className='p-0'></Command.List>
    </Command>
  );
});

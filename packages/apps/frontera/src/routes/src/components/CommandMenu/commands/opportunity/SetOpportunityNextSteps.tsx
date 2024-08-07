import { useState } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml';

export const SetOpportunityNextSteps = observer(() => {
  const store = useStore();
  const [value, setValue] = useState('');
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(
    (context.ids as string[])?.[0],
  );

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
          defaultHtmlValue={convertPlainTextToHtml(
            opportunity?.value?.nextSteps ?? '',
          )}
        />
      </CommandInput>

      <Command.List className='p-0'>
        <CommandItem
          leftAccessory={<InfoCircle />}
          className='data-[selected=true]:bg-white'
        >
          Use <code>⌘ + Enter</code> to save.
        </CommandItem>
      </Command.List>
    </Command>
  );
});

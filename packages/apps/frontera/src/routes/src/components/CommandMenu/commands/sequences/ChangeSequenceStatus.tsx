import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { FlowSequenceStatus } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeSequenceStatus = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = store.flowSequences.value.get(context.ids?.[0] as string);

  const handleSelect = (flowSequenceStatus: FlowSequenceStatus) => () => {
    if (!context.ids?.[0]) return;

    match(context.entity)
      .with('Sequence', () => {
        entity?.update((value) => {
          value.status = flowSequenceStatus;

          return value;
        });
      })
      .with('Sequences', () => {
        context.ids?.forEach((id) => {
          const sequence = store.flowSequences.value.get(id);

          sequence?.update((value) => {
            value.status = flowSequenceStatus;

            return value;
          });
        });
      })
      .otherwise(() => '');
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('SequenceCommands');
  };

  const label = match(context.entity)
    .with('Sequence', () => `Sequence - ${entity?.value?.name}`)
    .with('Sequences', () => `${context.ids?.length} sequences`)
    .otherwise(() => '');

  const status = entity?.value.status;

  return (
    <Command
      label='Change sequence status...'
      onKeyDown={(e) => {
        e.stopPropagation();
      }}
    >
      <CommandInput
        label={label}
        placeholder='Change sequence status...'
        onKeyDownCapture={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            store.ui.commandMenu.setOpen(false);
          }
        }}
      />

      <Command.List>
        <CommandItem
          key={FlowSequenceStatus.Active}
          onSelect={handleSelect(FlowSequenceStatus.Active)}
          rightAccessory={
            status === FlowSequenceStatus.Active ? <Check /> : null
          }
        >
          Active
        </CommandItem>
        <CommandItem
          key={FlowSequenceStatus.Paused}
          onSelect={handleSelect(FlowSequenceStatus.Paused)}
          rightAccessory={
            status === FlowSequenceStatus.Paused ? <Check /> : null
          }
        >
          Paused
        </CommandItem>
        <CommandItem
          key={FlowSequenceStatus.Inactive}
          onSelect={handleSelect(FlowSequenceStatus.Inactive)}
          rightAccessory={
            status === FlowSequenceStatus.Inactive ? <Check /> : null
          }
        >
          Inactive
        </CommandItem>
      </Command.List>
    </Command>
  );
});

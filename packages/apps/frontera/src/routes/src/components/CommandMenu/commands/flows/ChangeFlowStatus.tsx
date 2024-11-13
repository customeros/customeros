import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeFlowStatus = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = store.flows.value.get(context.ids?.[0] as string);

  const handleSelect = (flowStatus: FlowStatus) => () => {
    if (!context.ids?.[0]) return;

    match(context.entity)
      .with('Flow', () => {
        if (flowStatus === FlowStatus.On) {
          return store.ui.commandMenu.setType('StartFlow');
        }

        if (flowStatus === FlowStatus.Off) {
          return store.ui.commandMenu.setType('StopFlow');
        }
      })
      .with('Flows', () => {
        // todo refactor after bulk mutation is available
        context.ids?.forEach((id) => {
          const sequence = store.flows.value.get(id);

          sequence?.update((value) => {
            value.status = flowStatus;

            return value;
          });
        });
        store.ui.commandMenu.setOpen(false);
        store.ui.commandMenu.setType('FlowCommands');
      });
  };

  const label = match(context.entity)
    .with('Flow', () => `Flow - ${entity?.value?.name}`)
    .with('Flows', () => `${context.ids?.length} flows`)
    .otherwise(() => '');

  const status = entity?.value.status;

  return (
    <Command label='Change flow status...'>
      <CommandInput
        label={label}
        placeholder='Change flow status...'
        onKeyDownCapture={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            store.ui.commandMenu.setOpen(false);
          }
        }}
      />

      <Command.List>
        <CommandItem
          key={FlowStatus.On}
          onSelect={handleSelect(FlowStatus.On)}
          rightAccessory={status === FlowStatus.On ? <Check /> : null}
        >
          Live
        </CommandItem>
        <CommandItem
          key={FlowStatus.Off}
          onSelect={handleSelect(FlowStatus.Off)}
          rightAccessory={status === FlowStatus.Off ? <Check /> : null}
        >
          Stopped
        </CommandItem>
      </Command.List>
    </Command>
  );
});

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

    if (flowStatus === FlowStatus.Active) {
      store.ui.commandMenu.setType('StartFlow');

      return;
    }

    if (flowStatus === FlowStatus.Paused) {
      store.ui.commandMenu.setType('PauseFlow');

      return;
    }

    match(context.entity)
      .with('Flow', () => {
        entity?.update((value) => {
          value.status = flowStatus;

          return value;
        });
      })
      .with('Flows', () => {
        context.ids?.forEach((id) => {
          const sequence = store.flows.value.get(id);

          sequence?.update((value) => {
            value.status = flowStatus;

            return value;
          });
        });
      })
      .otherwise(() => '');
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('FlowCommands');
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
          key={FlowStatus.Active}
          onSelect={handleSelect(FlowStatus.Active)}
          rightAccessory={status === FlowStatus.Active ? <Check /> : null}
        >
          Live
        </CommandItem>
        <CommandItem
          key={FlowStatus.Paused}
          onSelect={handleSelect(FlowStatus.Paused)}
          rightAccessory={status === FlowStatus.Paused ? <Check /> : null}
        >
          Stopped
        </CommandItem>
        <CommandItem
          key={FlowStatus.Inactive}
          onSelect={handleSelect(FlowStatus.Inactive)}
          rightAccessory={status === FlowStatus.Inactive ? <Check /> : null}
        >
          Not Started
        </CommandItem>
      </Command.List>
    </Command>
  );
});

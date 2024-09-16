import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { FlowSequenceStatus } from '@graphql/types';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const ChangeSequenceStatusSubItemGroup = observer(() => {
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

  const status = entity?.value.status;

  return (
    <>
      <CommandSubItem
        icon={null}
        rightLabel='Live'
        leftLabel='Change status'
        key={FlowSequenceStatus.Active}
        onSelectAction={handleSelect(FlowSequenceStatus.Active)}
        rightAccessory={status === FlowSequenceStatus.Active ? <Check /> : null}
      />
      <CommandSubItem
        icon={null}
        rightLabel='Paused'
        leftLabel='Change status'
        key={FlowSequenceStatus.Paused}
        onSelectAction={handleSelect(FlowSequenceStatus.Paused)}
        rightAccessory={status === FlowSequenceStatus.Paused ? <Check /> : null}
      />
      <CommandSubItem
        icon={null}
        rightLabel='Not Started'
        leftLabel='Change status'
        key={FlowSequenceStatus.Inactive}
        onSelectAction={handleSelect(FlowSequenceStatus.Inactive)}
        rightAccessory={
          status === FlowSequenceStatus.Inactive ? <Check /> : null
        }
      />
    </>
  );
});

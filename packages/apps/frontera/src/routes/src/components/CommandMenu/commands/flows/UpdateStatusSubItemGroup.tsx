import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { Columns03 } from '@ui/media/icons/Columns03.tsx';

import { flowKeywords } from './keywords';

export const UpdateStatusSubItemGroup = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedIds = context.ids;

  const entity = store.flows.value.get(context.ids?.[0] as string);

  const isSelected = () => {
    if (selectedIds.length > 1) {
      return;
    } else {
      const flow = store.flows.value.get(selectedIds[0]);

      return flow?.value.status;
    }
  };

  const handleUpdateStatus = (status: FlowStatus) => {
    if (!context.ids?.[0]) return;

    if (status === FlowStatus.Active) {
      store.ui.commandMenu.setType('StartFlow');

      return;
    }

    if (status === FlowStatus.Paused) {
      store.ui.commandMenu.setType('PauseFlow');

      return;
    }

    match(context.entity)
      .with('Flow', () => {
        entity?.update((value) => {
          value.status = status;

          return value;
        });
      })
      .with('Flows', () => {
        context.ids?.forEach((id) => {
          const flow = store.flows.value.get(id);

          flow?.update((value) => {
            value.status = status;

            return value;
          });
        });
      })
      .otherwise(() => '');
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('FlowCommands');
  };

  return (
    <>
      <CommandSubItem
        rightLabel='Live'
        icon={<Columns03 />}
        leftLabel='Change flow status'
        keywords={[...flowKeywords.status_update, 'live']}
        rightAccessory={isSelected() === FlowStatus.Active ? <Check /> : null}
        onSelectAction={() => {
          handleUpdateStatus(FlowStatus.Active);
        }}
      />

      <CommandSubItem
        icon={<Columns03 />}
        rightLabel='Stopped'
        leftLabel='Change flow status'
        keywords={[...flowKeywords.status_update, 'stopped', 'paused']}
        rightAccessory={isSelected() === FlowStatus.Paused ? <Check /> : null}
        onSelectAction={() => {
          handleUpdateStatus(FlowStatus.Paused);
        }}
      />

      <CommandSubItem
        icon={<Columns03 />}
        rightLabel='Not started'
        leftLabel='Change flow status'
        keywords={[...flowKeywords.status_update, 'not started']}
        rightAccessory={isSelected() === FlowStatus.Inactive ? <Check /> : null}
        onSelectAction={() => {
          handleUpdateStatus(FlowStatus.Inactive);
        }}
      />
    </>
  );
});

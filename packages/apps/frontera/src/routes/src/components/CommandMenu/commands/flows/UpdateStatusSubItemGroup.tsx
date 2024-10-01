import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { Activity } from '@ui/media/icons/Activity';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

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
      store.ui.commandMenu.setType('StopFlow');

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
        icon={<Activity />}
        leftLabel='Change flow status'
        keywords={flowKeywords.status_update}
        rightAccessory={isSelected() === FlowStatus.Active ? <Check /> : null}
        onSelectAction={() => {
          handleUpdateStatus(FlowStatus.Active);
        }}
      />

      <CommandSubItem
        icon={<Activity />}
        rightLabel='Stopped'
        leftLabel='Change flow status'
        keywords={flowKeywords.status_update}
        rightAccessory={isSelected() === FlowStatus.Paused ? <Check /> : null}
        onSelectAction={() => {
          handleUpdateStatus(FlowStatus.Paused);
        }}
      />

      <CommandSubItem
        icon={<Activity />}
        rightLabel='Not started'
        leftLabel='Change flow status'
        keywords={flowKeywords.status_update}
        rightAccessory={isSelected() === FlowStatus.Inactive ? <Check /> : null}
        onSelectAction={() => {
          handleUpdateStatus(FlowStatus.Inactive);
        }}
      />
    </>
  );
});

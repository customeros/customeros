import { observer } from 'mobx-react-lite';
import { getStageFromColumn } from '@opportunities/components/ProspectsBoard/columns';

import { InternalStage } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

type OpportunityStage = InternalStage | string;

export const ChooseOpportunityStage = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const stages = store.tableViewDefs
    .getById(store.tableViewDefs.opportunitiesPreset ?? '')
    ?.value?.columns?.map((column) => ({
      value: getStageFromColumn(column),
      label: column.name,
    }))
    // TODO: remove this when the create mutation supports internal stages
    .filter(
      (o) =>
        ![InternalStage.ClosedLost, InternalStage.ClosedWon].includes(o.value),
    );

  const label = 'Choose stage';

  const handleSelect = (stage: OpportunityStage) => () => {
    store.ui.commandMenu.setType('ChooseOpportunityOrganization');
    store.ui.commandMenu.setContext({
      ...context,
      meta: { ...context.meta, stage },
    });
  };

  return (
    <Command
      onKeyDown={(e) => {
        e.stopPropagation();
      }}
    >
      <CommandInput label={label} placeholder='Choose stage...' />

      <Command.List>
        {stages?.map((option) => (
          <CommandItem key={option.value} onSelect={handleSelect(option.value)}>
            {option.label}
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});

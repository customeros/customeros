import { useMemo, useState, useEffect } from 'react';

import { useKeys } from 'rooks';
import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Currency, InternalStage } from '@graphql/types';

import { getColumns } from './columns';
import { PipelineMetrics } from '../PipelineMetrics';
import { KanbanColumn } from '../KanbanColumn/KanbanColumn';

export const ProspectsBoard = observer(() => {
  const store = useStore();
  const [focused, setFocused] = useState<string | null>(null);

  const opportunitiesPresetId = store.tableViewDefs.opportunitiesPreset;
  const viewDef = store.tableViewDefs.getById(opportunitiesPresetId ?? '');
  const stageLikelihoods = new Map(
    store.settings.tenant.value?.opportunityStages.map((s) => [
      s.value,
      s.likelihoodRate,
    ]),
  );

  const opportunities = store.opportunities.toComputedArray((arr) => {
    arr = arr.filter((opp) => opp.value.internalType === 'NBO');

    return arr;
  });

  const currency = store.settings.tenant.value?.baseCurrency || Currency.Usd;
  const count = opportunities.length;
  const totalArr = opportunities.reduce(
    (acc, card) => acc + card.value.maxAmount,
    0,
  );
  const totalWeightedArr = opportunities.reduce((acc, opp) => {
    const externalStage = opp.value.externalStage;
    const internalStage = opp.value.internalStage;
    const likelihoodRate = match(internalStage)
      .with(InternalStage.ClosedLost, () => 0)
      .with(InternalStage.ClosedWon, () => 100)
      .otherwise(() => stageLikelihoods.get(externalStage) ?? 0);

    return acc + opp.value.maxAmount * (likelihoodRate / 100);
  }, 0);

  const columns = useMemo(() => {
    return getColumns(viewDef?.value);
  }, [
    viewDef?.value.columns.reduce((acc, c) => acc + c.columnId + c.visible, ''),
  ]);

  const onDragEnd = (result: DropResult): void => {
    if (!result?.destination?.droppableId) return;

    const id = result.draggableId;
    const opportunity = store.opportunities.value.get(id);

    opportunity?.update((org) => {
      const destinationStage = result.destination?.droppableId;

      if (
        [
          InternalStage.Open,
          InternalStage.ClosedLost,
          InternalStage.ClosedWon,
        ].includes(destinationStage as InternalStage)
      ) {
        org.internalStage = destinationStage as InternalStage;
      } else {
        org.internalStage = InternalStage.Open;
        org.externalStage = destinationStage ?? 'STAGE1';
      }

      return org;
    });
  };

  const handleFocus = (id: string) => {
    setFocused(id);
    store.ui.commandMenu.setType('OpportunityCommands');
    store.ui.commandMenu.setContext({ entity: 'Opportunity', ids: [id] });
  };

  const handleBlur = () => {
    setFocused(null);
    store.ui.commandMenu.setType('OpportunityHub');
  };

  useEffect(() => {
    store.ui.commandMenu.setType('OpportunityHub');
  }, []);

  useEffect(() => {
    store.ui.setSearchCount(count);
  }, [count]);

  useKeys(
    ['Shift', 'S'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      store.ui.commandMenu.setType('ChangeStage');
      store.ui.commandMenu.setOpen(true);
    },
    { when: !!focused && !store.ui.commandMenu.isOpen },
  );

  useKeys(
    ['Shift', 'R'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      store.ui.commandMenu.setType('RenameOpportunityName');
      store.ui.commandMenu.setOpen(true);
    },
    { when: !!focused && !store.ui.commandMenu.isOpen },
  );
  useKeys(
    ['Shift', 'O'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      store.ui.commandMenu.setType('AssignOwner');
      store.ui.commandMenu.setOpen(true);
    },
    { when: !!focused && !store.ui.commandMenu.isOpen },
  );

  useModKey(
    'Backspace',
    () => {
      store.ui.commandMenu.setType('DeleteConfirmationModal');
      store.ui.commandMenu.setOpen(true);
    },
    { when: !!focused && !store.ui.commandMenu.isOpen },
  );

  return (
    <>
      <PipelineMetrics
        count={count}
        currency={currency}
        totalArr={totalArr}
        totalWeightedArr={totalWeightedArr}
      />

      <DragDropContext onDragEnd={onDragEnd}>
        <div className='flex flex-grow px-4 space-x-2'>
          {(columns ?? []).map((column) => {
            return (
              <KanbanColumn
                key={column.name}
                onBlur={handleBlur}
                focusedId={focused}
                stage={column.stage}
                onFocus={handleFocus}
                columnId={column.columnId}
                filterFns={column.filterFns ?? []}
                isLoading={store.organizations.isLoading}
              />
            );
          })}
          <div className='w-6'></div>
        </div>
      </DragDropContext>
    </>
  );
});

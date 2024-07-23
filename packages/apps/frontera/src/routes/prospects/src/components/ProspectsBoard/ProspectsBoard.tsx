import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { InternalStage } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';

import { getColumns } from './columns';
import { KanbanColumn } from '../KanbanColumn/KanbanColumn';

export const ProspectsBoard = observer(() => {
  const store = useStore();

  const opportunitiesPresetId = store.tableViewDefs.opportunitiesPreset;
  const viewDef = store.tableViewDefs.getById(opportunitiesPresetId ?? '');

  const columns = getColumns(viewDef?.value);

  const allOpportunities = store.opportunities.toComputedArray((arr) => {
    return arr.filter(
      (opp) =>
        (opp.value.internalStage === InternalStage.Open ||
          opp.value.internalStage === InternalStage.ClosedLost ||
          opp.value.internalStage === InternalStage.ClosedWon) &&
        opp.value.internalType === 'NBO',
    );
  });

  const onDragEnd = (result: DropResult): void => {
    if (!result.destination || !result.destination.droppableId) return;
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

  return (
    <>
      <div className='flex flex-col text-gray-700 overflow-auto'>
        <div className='px-4 mt-3'>
          <h1 className='text-xl font-bold'>Opportunities</h1>
        </div>

        <DragDropContext onDragEnd={onDragEnd}>
          <div className='flex flex-grow px-4 mt-4 space-x-2 h-[calc(100vh-10px)] overflow-y-scroll '>
            {(columns ?? [])
              ?.filter((p) => p.visible)
              .map((column) => {
                const items = allOpportunities.filter((opp) =>
                  column.filterFns?.reduce(
                    (acc, fn) => acc && fn(opp.value),
                    true,
                  ),
                );

                return (
                  <KanbanColumn
                    cards={items}
                    key={column.name}
                    type={column.stage}
                    cardCount={items.length}
                    columnId={column.columnId}
                    isLoading={store.organizations.isLoading}
                  />
                );
              })}
            <div className='flex-shrink-0 w-6'></div>
          </div>
        </DragDropContext>
      </div>
    </>
  );
});

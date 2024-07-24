import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { useStore } from '@shared/hooks/useStore';
import { InternalStage, TableViewType } from '@graphql/types';
import { ViewSettings } from '@shared/components/ViewSettings';

import { getColumns } from './columns';
import { KanbanColumn } from '../KanbanColumn/KanbanColumn';

export const ProspectsBoard = observer(() => {
  const store = useStore();

  const opportunitiesPresetId = store.tableViewDefs.opportunitiesPreset;
  const viewDef = store.tableViewDefs.getById(opportunitiesPresetId ?? '');

  const columns = useMemo(() => {
    return getColumns(viewDef?.value);
  }, [viewDef?.value.columns.reduce((acc, c) => acc + c.columnId, '')]);

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
        <div className='px-4 mt-3 flex justify-between'>
          <h1 className='text-xl font-bold'>Opportunities</h1>
          <ViewSettings type={TableViewType.Opportunities} />
        </div>

        <DragDropContext onDragEnd={onDragEnd}>
          <div className='flex flex-grow px-4 mt-4 space-x-2 h-[calc(100vh-10px)] overflow-y-scroll '>
            {(columns ?? []).map((column) => {
              return (
                <KanbanColumn
                  key={column.name}
                  type={column.stage}
                  columnId={column.columnId}
                  filterFns={column.filterFns ?? []}
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

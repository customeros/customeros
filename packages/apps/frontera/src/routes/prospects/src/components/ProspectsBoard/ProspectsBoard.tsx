import { useState, useEffect } from 'react';

import { reaction } from 'mobx';
import { useDidMount } from 'rooks';
import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { OrganizationStage } from '@graphql/types';

import { useOrganizationsKanbanData } from '../../hooks';
import { KanbanColumn } from '../KanbanColumn/KanbanColumn.tsx';

type ISortedKanbanColumns = {
  target: OrganizationStore[];
  engaged: OrganizationStore[];
  closed_won: OrganizationStore[];
  interested: OrganizationStore[];
};

type ISortedColumnKey = keyof ISortedKanbanColumns;
export const ProspectsBoard = observer(() => {
  useOrganizationsKanbanData({ sorting: [] });
  const { newBusiness } = useStore();
  useEffect(() => {
    newBusiness.bootstrap();
  }, []);

  const [sortedColumns, setSortedColumns] = useState<ISortedKanbanColumns>({
    target: [],
    engaged: [],
    closed_won: [],
    interested: [],
  });

  useDidMount(() => {
    sortKanbanValues();
  });

  useEffect(() => {
    const dispose = reaction(() => newBusiness.value.size, sortKanbanValues);

    return () => {
      dispose();
    };
  }, []);

  const sortKanbanValues = () => {
    const sortedKanbanValues = {
      ...sortedColumns,
    };
    newBusiness.value.forEach((org) => {
      if (
        org.value?.stage === OrganizationStage.ClosedWon &&
        sortedKanbanValues.closed_won.findIndex(
          (o) => o.value.metadata.id === org.value.metadata.id,
        ) === -1
      ) {
        sortedKanbanValues.closed_won.push(org);

        return;
      }
      if (
        org.value?.stage === OrganizationStage.Engaged &&
        sortedKanbanValues.engaged.findIndex(
          (o) => o.value.metadata.id === org.value.metadata.id,
        ) === -1
      ) {
        sortedKanbanValues.engaged.push(org);

        return;
      }

      if (
        org.value?.stage === OrganizationStage.Interested &&
        sortedColumns.interested.findIndex(
          (o) => o.value.metadata.id === org.value.metadata.id,
        ) === -1
      ) {
        sortedKanbanValues.interested.push(org);

        return;
      }

      if (
        org.value?.stage === OrganizationStage.Target &&
        sortedColumns.target.findIndex(
          (o) => o.value.metadata.id === org.value.metadata.id,
        ) === -1
      ) {
        sortedKanbanValues.target.push(org);

        return;
      }
    });

    sortedKanbanValues.target.sort(
      (a, b) =>
        new Date(b.value.metadata.created).getTime() -
        new Date(a.value.metadata.created).getTime(),
    );
    sortedKanbanValues.engaged.sort(
      (a, b) =>
        new Date(b.value.metadata.created).getTime() -
        new Date(a.value.metadata.created).getTime(),
    );
    sortedKanbanValues.interested.sort(
      (a, b) =>
        new Date(b.value.metadata.created).getTime() -
        new Date(a.value.metadata.created).getTime(),
    );
    sortedKanbanValues.closed_won.sort(
      (a, b) =>
        new Date(b.value.metadata.created).getTime() -
        new Date(a.value.metadata.created).getTime(),
    );

    setSortedColumns(sortedKanbanValues);
  };

  const onDragEnd = (result: DropResult): void => {
    if (!result.destination || !result.destination.droppableId) return;
    const currentColumnKey =
      result.source.droppableId.toLowerCase() as ISortedColumnKey;
    const destinationColumnKey =
      result.destination.droppableId.toLowerCase() as ISortedColumnKey;
    const item = sortedColumns[currentColumnKey]?.at(result.source.index);
    if (!item) return;
    const newValues = {
      ...sortedColumns,
    };

    newValues[currentColumnKey].splice(result.source.index, 1);
    newValues[destinationColumnKey].splice(result.destination.index, 0, item);

    setSortedColumns((prev) => ({
      ...prev,
      ...newValues,
    }));
    item.updateStage(result?.destination.droppableId as OrganizationStage);
  };

  const hasMorePages = newBusiness.page < newBusiness.totalPages;

  return (
    <>
      <div className='flex flex-col w-screen h-screen overflow-auto text-gray-700 '>
        <div className='px-4 mt-3'>
          <h1 className='text-xl font-bold'>New business</h1>
        </div>

        <DragDropContext onDragEnd={onDragEnd}>
          <div className='flex flex-grow px-4 mt-4 space-x-2 overflow-auto'>
            <KanbanColumn
              type={OrganizationStage.Target}
              title='Target'
              cardCount={sortedColumns.target.length}
              cards={sortedColumns.target}
              isLoading={newBusiness.isLoading}
            />
            <KanbanColumn
              type={OrganizationStage.Interested}
              title='Interested'
              cardCount={sortedColumns.interested.length}
              cards={sortedColumns.interested}
              isLoading={newBusiness.isLoading}
            />
            <KanbanColumn
              type={OrganizationStage.Engaged}
              title='Engaged'
              cardCount={sortedColumns.engaged.length}
              cards={sortedColumns.engaged}
              isLoading={newBusiness.isLoading}
            />

            <KanbanColumn
              type={OrganizationStage.ClosedWon}
              title='Closed Won'
              cardCount={sortedColumns.closed_won.length}
              cards={sortedColumns.closed_won}
              isLoading={newBusiness.isLoading}
            />
            <div className='flex-shrink-0 w-6'></div>
          </div>
        </DragDropContext>
        {hasMorePages && (
          <Button
            variant='ghost'
            colorScheme='primary'
            onClick={() => newBusiness?.loadMore()}
          >
            <Plus className='mr-2' />
            Load more
          </Button>
        )}
      </div>
    </>
  );
});

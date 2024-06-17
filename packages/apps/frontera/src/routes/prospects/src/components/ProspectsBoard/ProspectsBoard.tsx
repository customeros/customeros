import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { useStore } from '@shared/hooks/useStore';
import { Opportunity, InternalStage } from '@graphql/types';

import { KanbanColumn } from '../KanbanColumn/KanbanColumn.tsx';

export const ProspectsBoard = observer(() => {
  const store = useStore();

  const identified = store.opportunities.toComputedArray((arr) => {
    return arr.filter(
      (org) =>
        org.value.internalStage === InternalStage.Open &&
        org.value.externalStage === 'Identified',
    );
  });

  const committed = store.opportunities.toComputedArray((arr) => {
    return arr.filter(
      (org) =>
        org.value.internalStage === InternalStage.Open &&
        org.value.externalStage === 'Commited',
    );
  });

  const lost = store.opportunities.toComputedArray((arr) => {
    return arr.filter(
      (org) => org.value.internalStage === InternalStage.ClosedLost,
    );
  });

  const won = store.opportunities.toComputedArray((arr) => {
    return arr.filter(
      (org) => org.value.internalStage === InternalStage.ClosedWon,
    );
  });

  const onDragEnd = (result: DropResult): void => {
    if (!result.destination || !result.destination.droppableId) return;
    const id = result.draggableId;

    const opportunity = store.opportunities.value.get(id);

    if (
      result.destination.droppableId === 'Identified' ||
      result.destination.droppableId === 'Commited'
    ) {
      opportunity?.update((org) => {
        org.externalStage = result?.destination
          ?.droppableId as Opportunity['externalStage'];

        return org;
      });

      opportunity?.update(
        (org) => {
          org.internalStage = InternalStage.Open;

          return org;
        },
        { mutate: false },
      );
    } else {
      opportunity?.update((org) => {
        org.internalStage = result?.destination?.droppableId as InternalStage;

        return org;
      });
    }
  };

  return (
    <>
      <div className='flex flex-col w-screen h-[calc(100vh-10px)] text-gray-700 overflow-auto'>
        <div className='px-4 mt-3'>
          <h1 className='text-xl font-bold'>Opportunities</h1>
        </div>

        <DragDropContext onDragEnd={onDragEnd}>
          <div className='flex flex-grow px-4 mt-4 space-x-2 '>
            <KanbanColumn
              title='Identified'
              cards={identified}
              cardCount={identified.length}
              type={'Identified'}
              isLoading={store.organizations.isLoading}
              createOrganization={store.organizations.create}
            />

            <KanbanColumn
              title='Commited'
              cards={committed}
              cardCount={committed.length}
              type={'Commited'}
              isLoading={store.organizations.isLoading}
              createOrganization={store.organizations.create}
            />

            <KanbanColumn
              title='Lost'
              cards={lost}
              cardCount={lost.length}
              type={InternalStage.ClosedLost}
              isLoading={store.opportunities.isLoading}
              createOrganization={store.organizations.create}
            />

            <KanbanColumn
              title='Won'
              cards={won}
              cardCount={won.length}
              type={InternalStage.ClosedWon}
              isLoading={store.opportunities.isLoading}
              createOrganization={store.organizations.create}
            />
            <div className='flex-shrink-0 w-6'></div>
          </div>
        </DragDropContext>
      </div>
    </>
  );
});

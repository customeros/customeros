import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { useStore } from '@shared/hooks/useStore';
import { InternalStage, OrganizationStage } from '@graphql/types';

import { KanbanColumn } from '../KanbanColumn/KanbanColumn.tsx';

export const ProspectsBoard = observer(() => {
  const store = useStore();

  // const sortByCreatedAt = (a: Store<Organization>, b: Store<Organization>) =>
  //   new Date(b.value.metadata.created).getTime() -
  //   new Date(a.value.metadata.created).getTime();

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
    // const item = store.organizations.value.get(id);
    const opportunity = store.opportunities.value.get(id);

    if (
      result.destination.droppableId === 'Identified' ||
      result.destination.droppableId === 'Commited'
    ) {
      opportunity?.update((org) => {
        org.externalStage = result?.destination
          ?.droppableId as OrganizationStage;

        return org;
      });
    } else {
      opportunity?.update((org) => {
        org.internalStage = result?.destination?.droppableId as InternalStage;

        return org;
      });
    }
  };

  return (
    <>
      <div className='flex flex-col w-screen h-screen overflow-auto text-gray-700 '>
        <div className='px-4 mt-3'>
          <h1 className='text-xl font-bold'>Opportunities</h1>
        </div>

        <DragDropContext onDragEnd={onDragEnd}>
          <div className='flex flex-grow px-4 mt-4 space-x-2 overflow-auto'>
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

import { Store } from '@store/store';
import { observer } from 'mobx-react-lite';
import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { useStore } from '@shared/hooks/useStore';
import { Organization, OrganizationStage } from '@graphql/types';

import { KanbanColumn } from '../KanbanColumn/KanbanColumn.tsx';

export const ProspectsBoard = observer(() => {
  const store = useStore();

  const sortByCreatedAt = (a: Store<Organization>, b: Store<Organization>) =>
    new Date(b.value.metadata.created).getTime() -
    new Date(a.value.metadata.created).getTime();

  const engaged = store.organizations.toComputedArray((arr) => {
    return arr
      .filter((org) => org.value.stage === OrganizationStage.Engaged)
      .sort(sortByCreatedAt);
  });

  const trial = store.organizations.toComputedArray((arr) => {
    return arr
      .filter((org) => org.value.stage === OrganizationStage.Trial)
      .sort(sortByCreatedAt);
  });

  const readyToBuy = store.organizations.toComputedArray((arr) => {
    return arr
      .filter((org) => org.value.stage === OrganizationStage.ReadyToBuy)
      .sort(sortByCreatedAt);
  });

  const onboarding = store.organizations.toComputedArray((arr) => {
    return arr
      .filter((org) => org.value.stage === OrganizationStage.Onboarding)
      .sort(sortByCreatedAt);
  });

  const onDragEnd = (result: DropResult): void => {
    if (!result.destination || !result.destination.droppableId) return;
    const id = result.draggableId;
    const item = store.organizations.value.get(id);

    item?.update((org) => {
      org.stage = result?.destination?.droppableId as OrganizationStage;

      return org;
    });
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
              title='Engaged'
              cards={engaged}
              cardCount={engaged.length}
              type={OrganizationStage.Engaged}
              isLoading={store.organizations.isLoading}
              createOrganization={store.organizations.create}
            />

            <KanbanColumn
              title='Trial'
              cards={trial}
              cardCount={trial.length}
              type={OrganizationStage.Trial}
              isLoading={store.organizations.isLoading}
              createOrganization={store.organizations.create}
            />

            <KanbanColumn
              title='Ready to Buy'
              cards={readyToBuy}
              cardCount={readyToBuy.length}
              type={OrganizationStage.ReadyToBuy}
              isLoading={store.organizations.isLoading}
              createOrganization={store.organizations.create}
            />

            <KanbanColumn
              title='Won'
              cards={onboarding}
              cardCount={onboarding.length}
              type={OrganizationStage.Onboarding}
              isLoading={store.organizations.isLoading}
              createOrganization={store.organizations.create}
            />
            <div className='flex-shrink-0 w-6'></div>
          </div>
        </DragDropContext>
      </div>
    </>
  );
});

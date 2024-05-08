'use client';

import { useMemo, useCallback } from 'react';

import { DropResult, DragDropContext } from '@hello-pangea/dnd';

import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { Organization, OrganizationStage } from '@graphql/types';

import { KanbanColumn } from '../KanbanColumn/KanbanColumn';
import {
  useOrganizationsKanbanData,
  useOrganizationsPageMethods,
} from '../../hooks';

interface CategorizedOrganizations {
  lead: Organization[];
  target: Organization[];
  engaged: Organization[];
  nurture: Organization[];
  contracted: Organization[];
  interested: Organization[];
  uncategorized: Organization[];
}
export const ProspectsBoard = () => {
  const { data, isFetching, isLoading, hasNextPage, fetchNextPage } =
    useOrganizationsKanbanData({ sorting: [] });
  const { updateOrganization } = useOrganizationsPageMethods();

  function categorizeAndSortOrganizations(
    orgs: Organization[],
  ): CategorizedOrganizations {
    const categorized: CategorizedOrganizations = {
      uncategorized: [],
      contracted: [],
      engaged: [],
      interested: [],
      lead: [],
      nurture: [],
      target: [],
    };

    orgs.forEach((org) => {
      if (org.isCustomer) {
        return;
      }

      if (!org?.stage?.length) {
        categorized.uncategorized.push(org);

        return;
      }
      if (org?.stage === OrganizationStage.Contracted) {
        categorized.contracted.push(org);

        return;
      }
      if (org?.stage === OrganizationStage.Engaged) {
        categorized.engaged.push(org);

        return;
      }

      if (org?.stage === OrganizationStage.Interested) {
        categorized.interested.push(org);

        return;
      }
      if (org?.stage === OrganizationStage.Lead) {
        categorized.lead.push(org);

        return;
      }
      if (org?.stage === OrganizationStage.Target) {
        categorized.target.push(org);

        return;
      }
      if (org?.stage === OrganizationStage.Nurture) {
        categorized.nurture.push(org);

        return;
      }
    });

    // Sort each category by a date property, assuming `createdAt` exists on the organization type
    Object.keys(categorized).forEach((key) => {
      categorized[key as keyof CategorizedOrganizations].sort(
        (a, b) =>
          new Date(b.metadata.created).getTime() -
          new Date(a.metadata.created).getTime(),
      );
    });

    return categorized;
  }

  const categorized = useMemo(
    () => categorizeAndSortOrganizations(data || []),
    [data],
  );

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);
  const onDragEnd = (result: DropResult): void => {
    if (
      result.type === 'COLUMN' &&
      result.source.droppableId !== result.destination?.droppableId
    ) {
      updateOrganization.mutate({
        input: {
          id: result.draggableId,
          stage: result?.destination?.droppableId as OrganizationStage,
        },
      });
    }
  };

  return (
    <>
      <div className='flex flex-col w-screen h-screen overflow-auto text-gray-700 '>
        <div className='px-4 mt-3'>
          <h1 className='text-xl font-bold'>New business</h1>
        </div>

        <DragDropContext onDragEnd={onDragEnd}>
          <div className='flex flex-grow px-4 mt-4 space-x-2 overflow-auto'>
            <KanbanColumn
              type={OrganizationStage.Lead}
              title='Lead'
              cardCount={categorized.lead.length}
              cards={categorized.lead}
              isLoading={isLoading}
            />
            <KanbanColumn
              type={OrganizationStage.Target}
              title='Target'
              cardCount={categorized.target.length}
              cards={categorized.target}
              isLoading={isLoading}
            />
            <KanbanColumn
              type={OrganizationStage.Interested}
              title='Interested'
              cardCount={categorized.interested.length}
              cards={categorized.interested}
              isLoading={isLoading}
            />
            <KanbanColumn
              type={OrganizationStage.Engaged}
              title='Engaged'
              cardCount={categorized.engaged.length}
              cards={categorized.engaged}
              isLoading={isLoading}
            />

            <KanbanColumn
              type={OrganizationStage.Contracted}
              title='Closed Won'
              cardCount={categorized.contracted.length}
              cards={categorized.contracted}
              isLoading={isLoading}
            />
            <div className='flex-shrink-0 w-6'></div>
          </div>
        </DragDropContext>
        {hasNextPage && (
          <Button
            onClick={handleFetchMore}
            variant='ghost'
            colorScheme='primary'
          >
            <Plus className='mr-2' />
            Load more
          </Button>
        )}
      </div>
    </>
  );
};

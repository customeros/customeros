'use client';

import { useMemo, useCallback } from 'react';

import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { LogEntry, Organization, LastTouchpointType } from '@graphql/types';

import { useOrganizationsKanbanData } from '../../hooks';
import { KanbanColumn } from '../KanbanColumn/KanbanColumn';

interface CategorizedOrganizations {
  new: Organization[];
  contacted: Organization[];
  abandoned: Organization[];
  opportunity: Organization[];
}
export const ProspectsBoard = () => {
  const { data, isFetching, isLoading, hasNextPage, fetchNextPage } =
    useOrganizationsKanbanData({ sorting: [] });

  function categorizeAndSortOrganizations(
    orgs: Organization[],
  ): CategorizedOrganizations {
    const categorized: CategorizedOrganizations = {
      new: [],
      contacted: [],
      opportunity: [],
      abandoned: [],
    };

    orgs.forEach((org) => {
      if (org.isCustomer) {
        return;
      }
      const tags = (
        org.lastTouchpoint?.lastTouchPointTimelineEvent as LogEntry
      )?.tags?.map((e) => e.name?.toLowerCase());

      if (
        org?.lastTouchpoint &&
        tags?.length > 0 &&
        ['abandoned', 'lost', 'killed'].some((event) => tags.includes(event))
      ) {
        categorized.abandoned.push(org);

        return;
      }
      if (org?.contracts && org?.contracts?.length > 0) {
        categorized.opportunity.push(org);

        return;
      }

      if (
        org?.lastTouchpoint &&
        org.lastTouchpoint?.lastTouchPointType ===
          LastTouchpointType.ActionCreated
      ) {
        categorized.new.push(org);

        return;
      }

      if (
        org.lastTouchpoint &&
        org.lastTouchpoint.lastTouchPointTimelineEvent
      ) {
        categorized.contacted.push(org);

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

  return (
    <>
      <div className='flex flex-col w-screen h-screen overflow-auto text-gray-700 '>
        <div className='px-10 mt-3'>
          <h1 className='text-2xl font-bold'>Prospects</h1>
        </div>
        <div className='flex flex-grow px-10 mt-4 space-x-6 overflow-auto'>
          <KanbanColumn
            title='New prospects'
            cardCount={categorized.new.length}
            cards={categorized.new}
            isLoading={isLoading}
          />
          <KanbanColumn
            title='Contacted'
            cardCount={categorized.contacted.length}
            cards={categorized.contacted}
            isLoading={isLoading}
          />
          <KanbanColumn
            title='Opportunity'
            cardCount={categorized.opportunity.length}
            cards={categorized.opportunity}
            isLoading={isLoading}
          />
          <KanbanColumn
            title='Abandoned'
            cardCount={categorized.abandoned.length}
            cards={categorized.abandoned}
            isLoading={isLoading}
          />
          <div className='flex-shrink-0 w-6'></div>
        </div>
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

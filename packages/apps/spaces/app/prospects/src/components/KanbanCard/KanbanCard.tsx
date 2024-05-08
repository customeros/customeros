import React, { forwardRef } from 'react';
import { useRouter } from 'next/navigation';

import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Clock } from '@ui/media/icons/Clock';
import { Organization } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { Mail01 } from '@ui/media/icons/Mail01';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building06 } from '@ui/media/icons/Building06';
import { Building05 } from '@ui/media/icons/Building05';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

interface DraggableKanbanCardProps {
  index: number;
  card: Organization;
  noPointerEvents?: boolean;
}

export const DraggableKanbanCard = forwardRef<
  HTMLDivElement,
  DraggableKanbanCardProps
>(({ card, index, noPointerEvents }, _ref) => {
  return (
    <Draggable index={index} draggableId={card?.metadata.id}>
      {(provided, snapshot) => {
        return (
          <KanbanCard
            card={card}
            provided={provided}
            snapshot={snapshot}
            noPointerEvents={noPointerEvents}
          />
        );
      }}
    </Draggable>
  );
});

interface KanbanCardProps {
  card: Organization;
  noPointerEvents?: boolean;
  provided?: DraggableProvided;
  snapshot?: DraggableStateSnapshot;
}

export const KanbanCard: React.FC<KanbanCardProps> = ({
  card,
  provided,
  snapshot,
  noPointerEvents,
}) => {
  const router = useRouter();
  const ownerName = `${card?.owner?.firstName ? card?.owner?.firstName : ''}${
    card?.owner?.lastName && card?.owner?.firstName ? ' ' : ''
  }${card?.owner?.lastName ? card?.owner?.lastName : ''}`;

  return (
    <div
      className={cn(
        'cursor-grab relative flex flex-col items-start p-3 mt-3 bg-white rounded-lg cursor-pointer bg-opacity-90 group hover:bg-opacity-100 border border-gray-200 shadow-xs',
        {
          'shadow-lg': snapshot?.isDragging,
          'pointer-events-none': noPointerEvents,
        },
      )}
      ref={provided?.innerRef}
      {...provided?.draggableProps}
    >
      <div
        className='flex justify-between w-full'
        {...provided?.dragHandleProps}
      >
        <div className='flex '>
          {card.logo && (
            <Avatar
              name={`${card.name}`}
              size='xs'
              icon={<Building06 className='text-primary-500 size-3' />}
              className={cn(
                card.logo
                  ? 'mr-1 h-5 w-5'
                  : 'border border-primary-200 mr-1 h-5 w-5',
              )}
              src={card.logo || undefined}
              variant='roundedSquareSmall'
            />
          )}

          <Button
            size='sm'
            variant='link'
            className='text-sm font-medium shadow-none p-0'
            onClick={() => router.push(`/organization/${card.metadata.id}`)}
          >
            {card.name}
          </Button>
        </div>

        {card.timelineEventsTotalCount > 0 && (
          <Tooltip
            label={`Last touch point at ${DateTimeUtils.format(
              card.lastTouchpoint?.lastTouchPointAt,
              DateTimeUtils.dateWithAbreviatedMonth,
            )}`}
          >
            <div className='flex items-center'>
              <span className='mr-1 text-xs'>
                {card.timelineEventsTotalCount}
              </span>
              <Mail01 className='text-gray-500 size-3' />
            </div>
          </Tooltip>
        )}
      </div>
      <div className='flex mt-2 items-center'>
        <Building05 className='mr-1 text-gray-500' />
        <span className='text-sm text-gray-600'>
          {card?.employees} employees
        </span>
      </div>
      {(card.owner?.firstName || card.owner?.lastName) && (
        <div className='flex items-center mt-2'>
          <Tooltip label={`${card.owner.firstName} ${card.owner.lastName}`}>
            <Avatar
              name={ownerName}
              size='xs'
              icon={<User01 className='text-primary-500 size-3' />}
              className={cn(
                card.owner.profilePhotoUrl
                  ? 'mr-1 h-5 w-5'
                  : 'border border-primary-200 mr-1 h-5 w-5',
              )}
              src={card.owner.profilePhotoUrl || undefined}
            />
          </Tooltip>
          <div className='text-sm text-gray-500'>
            {card.owner.firstName} {card.owner.lastName}
          </div>
        </div>
      )}
      <div className='flex justify-between items-end w-full mt-2'>
        <div className='text-sm'>
          {/* todo use estimated opportunity size*/}
          {/*{card.accountDetails?.renewalSummary?.maxArrForecast &&*/}
          {/*  formatCurrency(2000, 2, 'USD')}*/}
          0$
        </div>

        <div className='text-xs font-medium text-gray-400 flex items-center'>
          <Clock className='mr-1' />
          {DateTimeUtils.timeAgo(card.metadata.lastUpdated, {
            addSuffix: true,
          })}
        </div>
      </div>
    </div>
  );
};

import React from 'react';
import { useRouter } from 'next/navigation';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Organization } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface KanbanCardProps {
  cardStyle: string;
  card: Organization;
}

export const KanbanCard: React.FC<KanbanCardProps> = ({ card, cardStyle }) => {
  const router = useRouter();

  return (
    <div
      className={cn(
        'relative flex flex-col items-start p-4 mt-3 bg-white rounded-lg cursor-pointer bg-opacity-90 group hover:bg-opacity-100 border-gray-100',
        cardStyle,
      )}
    >
      <Button
        size='sm'
        variant='link'
        className='text-sm font-medium shadow-none p-0'
        onClick={() => router.push(`/organization/${card.metadata.id}`)}
      >
        {card.name}
      </Button>
      <div className='flex items-baseline justify-between w-full mt-3 text-xs font-medium text-gray-400'>
        <span className='ml-1 leading-none'>
          Created{' '}
          {DateTimeUtils.timeAgo(card.metadata.created, { addSuffix: true })}
        </span>
        {card.owner?.firstName && (
          <Tooltip
            label={`Owner: ${card.owner.firstName} ${card.owner.lastName}`}
          >
            <Avatar
              name={''}
              size='xs'
              icon={<User01 className='text-primary-500 size-3' />}
              className={cn(
                card.owner.profilePhotoUrl
                  ? 'mr-1'
                  : 'border border-primary-200 mr-1',
              )}
              src={card.owner.profilePhotoUrl || undefined}
            />
          </Tooltip>
        )}
      </div>
    </div>
  );
};

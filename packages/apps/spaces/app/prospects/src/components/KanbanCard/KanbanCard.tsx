import React, { forwardRef } from 'react';
import { useRouter } from 'next/navigation';

import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building06 } from '@ui/media/icons/Building06';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Organization, OrganizationStage } from '@graphql/types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { useOrganizationsPageMethods } from '../../hooks';

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
  const { updateOrganization } = useOrganizationsPageMethods();

  const ownerName = `${card?.owner?.firstName ? card?.owner?.firstName : ''}${
    card?.owner?.lastName && card?.owner?.firstName ? ' ' : ''
  }${card?.owner?.lastName ? card?.owner?.lastName : ''}`;

  const handleChangeStage = (stage: OrganizationStage): void => {
    updateOrganization.mutate({
      input: {
        id: card.metadata.id,
        stage,
      },
    });
  };

  return (
    <div
      tabIndex={0}
      className={cn(
        'cursor-grab relative flex flex-col items-start p-2 pl-3 mb-2 bg-white rounded-lg border border-gray-200 shadow-xs hover:shadow-lg focus:border-primary-500 transition-all duration-200 ease-in-out',
        {
          'shadow-lg rotate-3': snapshot?.isDragging,
          'pointer-events-none': noPointerEvents,
        },
      )}
      ref={provided?.innerRef}
      {...provided?.draggableProps}
      {...provided?.dragHandleProps}
    >
      <div className='flex justify-between w-full items-center'>
        {card.logo && (
          <Avatar
            name={`${card.name}`}
            size='xs'
            icon={<Building06 className='text-primary-500 size-3' />}
            className={cn(
              card.logo ? 'h-5 w-5' : 'border border-primary-200  h-5 w-5',
            )}
            src={card.logo || undefined}
            variant='roundedSquareSmall'
          />
        )}

        <Button
          size='sm'
          variant='link'
          className='text-sm font-medium shadow-none p-0 no-underline'
          onClick={() => router.push(`/organization/${card.metadata.id}`)}
        >
          {card.name}
        </Button>
        <div className='flex items-center '>
          <Menu>
            <MenuButton aria-label='Stage' className='flex items-center mr-1'>
              <DotsVertical className='text-gray-500 w-4' />
            </MenuButton>
            <MenuList
              align='start'
              side='bottom'
              className='w-[200px] shadow-xl'
            >
              <MenuItem
                color='gray.700'
                onClick={() => handleChangeStage(OrganizationStage.Nurture)}
              >
                <HeartHand className='text-gray-500 mr-2' />
                Nurture
              </MenuItem>

              <MenuItem
                color='gray.700'
                onClick={() => handleChangeStage(OrganizationStage.ClosedLost)}
              >
                <BrokenHeart className='text-gray-500 mr-2' />
                Closed lost
              </MenuItem>
            </MenuList>
          </Menu>
          {(card.owner?.firstName || card.owner?.lastName) && (
            <Tooltip label={`${card.owner.firstName} ${card.owner.lastName}`}>
              <Avatar
                name={ownerName}
                textSizes={'xs'}
                size='xs'
                icon={<User01 className='text-primary-500 size-3' />}
                className={cn(
                  card.owner.profilePhotoUrl
                    ? ''
                    : 'border border-primary-200 text-xs',
                )}
                src={card.owner.profilePhotoUrl || undefined}
              />
            </Tooltip>
          )}
        </div>
      </div>
    </div>
  );
};

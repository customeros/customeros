import { useNavigate } from 'react-router-dom';
import React, { useState, forwardRef } from 'react';

import { Store } from '@store/store';
import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { User01 } from '@ui/media/icons/User01';
import { UserX01 } from '@ui/media/icons/UserX01.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Building06 } from '@ui/media/icons/Building06';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Organization, OrganizationStage } from '@graphql/types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface DraggableKanbanCardProps {
  index: number;
  card: Store<Organization>;
  noPointerEvents?: boolean;
}

export const DraggableKanbanCard = forwardRef<
  HTMLDivElement,
  DraggableKanbanCardProps
>(({ card, index, noPointerEvents }, _ref) => {
  return (
    <Draggable index={index} draggableId={card?.value?.metadata.id}>
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
  card: Store<Organization>;
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
  const navigate = useNavigate();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const ownerName = `${
    card?.value?.owner?.firstName ? card?.value?.owner?.firstName : ''
  }${card?.value?.owner?.lastName && card?.value?.owner?.firstName ? ' ' : ''}${
    card?.value?.owner?.lastName ? card?.value?.owner?.lastName : ''
  }`;

  const handleChangeStage = (stage: OrganizationStage): void => {
    card.update((org) => {
      org.stage = stage;

      return org;
    });
  };

  return (
    <div
      tabIndex={0}
      ref={provided?.innerRef}
      onMouseUp={() => {
        if (isMenuOpen) return;
        navigate(`/organization/${card.value?.metadata.id}`);
      }}
      {...provided?.draggableProps}
      {...provided?.dragHandleProps}
      className={cn(
        ' group/kanbanCard !cursor-pointer relative flex flex-col items-start p-2 pl-3 mb-2 bg-white rounded-lg border border-gray-200 shadow-xs hover:shadow-lg focus:border-primary-500 transition-all duration-200 ease-in-out',
        {
          '!shadow-lg cursor-grabbing': snapshot?.isDragging,
          'pointer-events-none': noPointerEvents,
        },
      )}
    >
      <div className='flex justify-between w-full items-center'>
        {card.value?.logo && (
          <Avatar
            name={`${card.value?.name}`}
            size='xs'
            icon={<Building06 className='text-primary-500 size-3' />}
            className={cn(
              card.value?.logo
                ? 'h-5 w-5'
                : 'border border-primary-200  h-5 w-5',
            )}
            src={card.value?.logo || undefined}
            variant='roundedSquareSmall'
          />
        )}
        <span
          role='navigation'
          className='text-sm font-medium shadow-none p-0 no-underline hover:no-underline focus:no-underline '
        >
          {card.value?.name}
        </span>

        <div className='flex items-center '>
          <Menu
            defaultOpen={false}
            open={isMenuOpen}
            onOpenChange={(status) => setIsMenuOpen(status)}
          >
            <MenuButton
              aria-label='Stage'
              className={
                'flex items-center mr-1 opacity-0 group-hover/kanbanCard:opacity-100 aria-[expanded=true]:opacity-100'
              }
            >
              <DotsVertical className='text-gray-500 w-4' />
            </MenuButton>
            <MenuList
              align='start'
              side='bottom'
              className='w-[200px] shadow-xl'
            >
              <MenuItem
                color='gray.700'
                onClick={() => {
                  handleChangeStage(OrganizationStage.Unqualified);
                }}
              >
                <UserX01 className='text-gray-500 mr-2' />
                Unqualify
              </MenuItem>

              <MenuItem
                color='gray.700'
                onClick={() => {
                  handleChangeStage(OrganizationStage.Lead);
                }}
              >
                <BrokenHeart className='text-gray-500 mr-2' />
                Make Lead
              </MenuItem>
              <MenuItem
                color='gray.700'
                onClick={() => {
                  handleChangeStage(OrganizationStage.Target);
                }}
              >
                <HeartHand className='text-gray-500 mr-2' />
                Nurture
              </MenuItem>
            </MenuList>
          </Menu>
          {(card.value?.owner?.firstName || card.value?.owner?.lastName) && (
            <Tooltip
              label={`${card.value?.owner.firstName} ${card.value?.owner.lastName}`}
            >
              <Avatar
                name={ownerName}
                textSize='xs'
                size='xs'
                icon={<User01 className='text-primary-500 size-3' />}
                className={cn(
                  card.value?.owner.profilePhotoUrl
                    ? ''
                    : 'border border-primary-200 text-xs',
                )}
                src={card.value?.owner.profilePhotoUrl || undefined}
              />
            </Tooltip>
          )}
        </div>
      </div>
    </div>
  );
};

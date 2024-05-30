import { useState } from 'react';

import { Store } from '@store/store';
import {
  Droppable,
  DroppableProvided,
  DroppableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { ResizableInput } from '@ui/form/Input';
import { Skeleton } from '@ui/feedback/Skeleton';
import { IconButton } from '@ui/form/IconButton';
import { Organization, OrganizationStage } from '@graphql/types';

import { KanbanCard, DraggableKanbanCard } from '../KanbanCard/KanbanCard';

interface CardColumnProps {
  title: string;
  cardCount: number;
  isLoading: boolean;
  cards: Store<Organization>[];
  createOrganization: () => void;
  type: OrganizationStage | 'new';
}

export const KanbanColumn = ({
  title,
  cardCount,
  cards,
  isLoading,
  type,
  createOrganization,
}: CardColumnProps) => {
  const [newData, setNewData] = useState<Array<{ id: string; name: string }>>(
    [],
  );
  const handleUpdateNewData = (id: string, newName: string) => {
    setNewData((prev) =>
      prev.map((item) => (item.id === id ? { ...item, name: newName } : item)),
    );
  };

  const handleRemoveNewData = (id: string) => {
    setNewData((prev) => prev.filter((item) => item.id !== id));
  };

  const handleSaveNewData = () => {
    createOrganization();
  };

  return (
    <div className='flex flex-col flex-shrink-0 w-72 bg-gray-100 rounded'>
      <div className='flex items-center justify-between flex-shrink-0 h-10 p-3 pb-0'>
        <div className='flex'>
          <span className='block text-sm font-semibold'>{title}</span>
          <span
            className={cn(
              'flex items-center justify-center w-5 h-5 ml-1 text-sm font-semibold  rounded text-gray-500 whitespace-nowrap',
            )}
          >
            • {cardCount}
          </span>
        </div>

        {title === 'New prospects' && (
          <IconButton
            aria-label={'Add new prospect'}
            icon={<Plus />}
            variant='ghost'
            size='xs'
            onClick={createOrganization}
          />
        )}
      </div>

      <Droppable
        droppableId={type}
        type={`COLUMN`}
        key={`kanban-columns-${title}`}
        renderClone={(provided, snapshot, rubric) => {
          return (
            <KanbanCard
              provided={provided}
              snapshot={snapshot}
              card={cards[rubric.source.index]}
            />
          );
        }}
      >
        {(
          dropProvided: DroppableProvided,
          dropSnapshot: DroppableStateSnapshot,
        ) => (
          <div
            className={cn('flex flex-col pb-2 overflow-auto p-3 min-h-[100%]', {
              'bg-gray-100': dropSnapshot?.isDraggingOver,
            })}
            ref={dropProvided.innerRef}
            {...dropProvided.droppableProps}
          >
            {newData.map((data) => (
              <div
                key={data.id}
                className={cn(
                  'relative flex flex-col items-start p-4 mt-3 bg-white rounded-lg cursor-pointer bg-opacity-90 group hover:bg-opacity-100',
                )}
              >
                <ResizableInput
                  value={data.name}
                  className='text-sm font-medium shadow-none p-0 min-h-5'
                  autoFocus
                  onChange={(e) => handleUpdateNewData(data.id, e.target.value)}
                />

                <div className='flex justify-end w-full'>
                  <IconButton
                    variant='ghost'
                    size='xs'
                    aria-label='Cancel'
                    className='p-1'
                    icon={<X />}
                    onClick={() => handleRemoveNewData(data.id)}
                  />
                  <IconButton
                    variant='ghost'
                    size='xs'
                    aria-label='Save'
                    className='p-1'
                    icon={<Check />}
                    onClick={handleSaveNewData}
                  />
                </div>
              </div>
            ))}

            {cards.map((card, index) => (
              <>
                <DraggableKanbanCard
                  index={index}
                  card={card}
                  noPointerEvents={dropSnapshot.isDraggingOver}
                  key={`card-${card.value.name}-${card.value.metadata.id}-${index}`}
                />
              </>
            ))}
            {isLoading && (
              <>
                <Skeleton className='h-[38px] min-h-[38px] rounded-lg mt-3' />
                <Skeleton className='h-[38px] min-h-[38px] rounded-lg mt-3' />
                <Skeleton className='h-[38px] min-h-[38px] rounded-lg mt-3' />
              </>
            )}
            {dropProvided.placeholder}
          </div>
        )}
      </Droppable>
    </div>
  );
};

import { useState } from 'react';

import { Store } from '@store/store';
import { observer } from 'mobx-react-lite';
import {
  Droppable,
  DroppableProvided,
  DroppableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { Skeleton } from '@ui/feedback/Skeleton';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Input, ResizableInput } from '@ui/form/Input';
import { Opportunity, InternalStage } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

import { KanbanCard, DraggableKanbanCard } from '../KanbanCard/KanbanCard';

interface CardColumnProps {
  columnId: number;
  cardCount: number;
  isLoading: boolean;
  cards: Store<Opportunity>[];
  type: string | InternalStage.ClosedLost | InternalStage.ClosedWon;
}

export const KanbanColumn = observer(
  ({ cards, type, isLoading, cardCount, columnId }: CardColumnProps) => {
    const store = useStore();
    const viewDef = store.tableViewDefs.getById(
      store.tableViewDefs.opportunitiesPreset ?? '',
    );
    const column = viewDef?.value.columns.find((c) => c.columnId === columnId);
    const [newData, setNewData] = useState<Array<{ id: string; name: string }>>(
      [],
    );

    const totalSum = formatCurrency(
      cards.reduce((acc, card) => acc + card.value.maxAmount, 0),
      0,
      store.settings.tenant.value?.baseCurrency as string,
    );

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      viewDef?.setColumnName(columnId, e.target.value);
    };

    const handleUpdateNewData = (id: string, newName: string) => {
      setNewData((prev) =>
        prev.map((item) =>
          item.id === id ? { ...item, name: newName } : item,
        ),
      );
    };

    const handleRemoveNewData = (id: string) => {
      setNewData((prev) => prev.filter((item) => item.id !== id));
    };

    const handleSaveNewData = () => {
      store.organizations.create();
    };

    return (
      <div className='flex flex-col flex-shrink-0 w-72 bg-gray-100 rounded h-full'>
        <div className='flex items-center justify-between p-3 pb-0'>
          <div className='flex flex-col items-center mb-2'>
            <div>
              <Input
                size='sm'
                variant='unstyled'
                value={column?.name}
                className='h-auto font-semibold min-h-[unset]'
                onChange={handleNameChange}
              />
            </div>
            <span className={cn('w-full text-sm font-medium text-gray-500')}>
              {`${totalSum} • ${cardCount}`}
            </span>
          </div>

          {column?.name === 'New prospects' && (
            <IconButton
              aria-label={'Add new prospect'}
              icon={<Plus />}
              variant='ghost'
              size='xs'
              onClick={handleSaveNewData}
            />
          )}
        </div>

        <Droppable
          droppableId={type}
          type={`COLUMN`}
          key={`kanban-columns-${columnId}`}
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
              className={cn('flex flex-col pb-2 overflow-auto p-3 ', {
                'bg-gray-100': dropSnapshot?.isDraggingOver,
              })}
              ref={dropProvided.innerRef}
              {...dropProvided.droppableProps}
              style={{ height: 'calc(100% - 40px)' }}
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
                    onChange={(e) =>
                      handleUpdateNewData(data.id, e.target.value)
                    }
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
                    card={card}
                    index={index}
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
  },
);

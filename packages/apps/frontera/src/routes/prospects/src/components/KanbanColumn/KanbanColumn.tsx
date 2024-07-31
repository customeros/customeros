import { useState } from 'react';

import { match } from 'ts-pattern';
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
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Percent03 } from '@ui/media/icons/Percent03';
import { Input, ResizableInput } from '@ui/form/Input';
import { Opportunity, InternalStage } from '@graphql/types';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { WinProbabilityModal } from './WinProbabilityModal';
import { KanbanCard, DraggableKanbanCard } from '../KanbanCard/KanbanCard';

interface CardColumnProps {
  columnId: number;
  isLoading: boolean;
  onBlur: () => void;
  focusedId: string | null;
  onFocus: (id: string) => void;
  filterFns: Array<(opportunity: Opportunity) => boolean>;
  stage: string | InternalStage.ClosedLost | InternalStage.ClosedWon;
}

export const KanbanColumn = observer(
  ({
    stage,
    onBlur,
    onFocus,
    columnId,
    focusedId,
    filterFns,
    isLoading,
  }: CardColumnProps) => {
    const { open, onOpen, onToggle } = useDisclosure();
    const store = useStore();
    const viewDef = store.tableViewDefs.getById(
      store.tableViewDefs.opportunitiesPreset ?? '',
    );
    const column = viewDef?.value.columns.find((c) => c.columnId === columnId);
    const [newData, setNewData] = useState<Array<{ id: string; name: string }>>(
      [],
    );

    const cards = store.opportunities.toComputedArray((arr) => {
      return arr.filter(
        (opp) =>
          opp.value.internalType === 'NBO' &&
          filterFns.every((fn) => fn(opp.value)),
      );
    });

    const totalSum = formatCurrency(
      cards.reduce((acc, card) => acc + card.value.maxAmount, 0),
      0,
      store.settings.tenant.value?.baseCurrency as string,
    );

    const stageLikelihoodRate = match(stage)
      .with(InternalStage.ClosedLost, () => 0)
      .with(InternalStage.ClosedWon, () => 100)
      .otherwise(() => {
        return (
          store.settings.tenant.value?.opportunityStages.find(
            (s) => s.value === stage,
          )?.likelihoodRate ?? 0
        );
      });

    const canEdit = match(stage)
      .with(InternalStage.ClosedLost, InternalStage.ClosedWon, () => false)
      .otherwise(() => true);

    const handleUpdateAllProbabilities = () => {
      cards.forEach((card) => {
        card.update((value) => {
          value.likelihoodRate = stageLikelihoodRate;
          value.amount = value.maxAmount * (stageLikelihoodRate / 100);

          return value;
        });
      });
    };

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      viewDef?.setColumnName(columnId, e.target.value);
    };

    const handleNameBlur = () => {
      viewDef?.save();
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
          <div className='flex flex-col items-center mb-2 w-full'>
            <div className='flex justify-between w-full'>
              <Tooltip
                asChild
                align='start'
                label={!canEdit ? 'This stage can’t be edited' : undefined}
              >
                <Input
                  size='sm'
                  variant='unstyled'
                  disabled={!canEdit}
                  value={column?.name}
                  onBlur={handleNameBlur}
                  onChange={handleNameChange}
                  className={cn(
                    'h-auto font-semibold min-h-[unset]',
                    !canEdit && 'cursor-not-allowed',
                  )}
                />
              </Tooltip>
              <Menu>
                <MenuButton asChild>
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    isDisabled={!canEdit}
                    icon={<DotsVertical />}
                    aria-label='Column options'
                  />
                </MenuButton>
                <MenuList>
                  <MenuItem onClick={onOpen}>
                    <Percent03 />
                    <span>Set win probability</span>
                  </MenuItem>
                </MenuList>
              </Menu>
            </div>
            <span className={cn('w-full text-sm font-medium text-gray-500')}>
              {`${totalSum} • ${cards.length}`}
            </span>
          </div>

          {column?.name === 'New prospects' && (
            <IconButton
              size='xs'
              icon={<Plus />}
              variant='ghost'
              onClick={handleSaveNewData}
              aria-label={'Add new prospect'}
            />
          )}
        </div>

        <Droppable
          type={`COLUMN`}
          droppableId={stage}
          key={`kanban-columns-${columnId}`}
          renderClone={(provided, snapshot, rubric) => {
            return (
              <KanbanCard
                onBlur={onBlur}
                onFocus={onFocus}
                provided={provided}
                snapshot={snapshot}
                card={cards[rubric.source.index]}
                isFocused={cards[rubric.source.index]?.id === focusedId}
              />
            );
          }}
        >
          {(
            dropProvided: DroppableProvided,
            dropSnapshot: DroppableStateSnapshot,
          ) => (
            <div
              ref={dropProvided.innerRef}
              className={cn('flex flex-col pb-2 overflow-auto p-3 ', {
                'bg-gray-100': dropSnapshot?.isDraggingOver,
              })}
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
                    autoFocus
                    value={data.name}
                    className='text-sm font-medium shadow-none p-0 min-h-5'
                    onChange={(e) =>
                      handleUpdateNewData(data.id, e.target.value)
                    }
                  />

                  <div className='flex justify-end w-full'>
                    <IconButton
                      size='xs'
                      icon={<X />}
                      variant='ghost'
                      className='p-1'
                      aria-label='Cancel'
                      onClick={() => handleRemoveNewData(data.id)}
                    />
                    <IconButton
                      size='xs'
                      variant='ghost'
                      className='p-1'
                      icon={<Check />}
                      aria-label='Save'
                      onClick={handleSaveNewData}
                    />
                  </div>
                </div>
              ))}

              {cards.map((card, index) => (
                <DraggableKanbanCard
                  card={card}
                  index={index}
                  onBlur={onBlur}
                  onFocus={onFocus}
                  isFocused={card.id === focusedId}
                  noPointerEvents={dropSnapshot.isDraggingOver}
                  key={`card-${card.value.name}-${card.value.metadata.id}-${index}`}
                />
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

        <WinProbabilityModal
          open={open}
          stage={stage}
          onToggle={onToggle}
          columnName={column?.name ?? ''}
          onUpdateProbability={handleUpdateAllProbabilities}
        />
      </div>
    );
  },
);

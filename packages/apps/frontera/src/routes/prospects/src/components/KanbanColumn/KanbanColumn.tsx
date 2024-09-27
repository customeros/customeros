import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import {
  Droppable,
  DroppableProvided,
  DroppableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Plus } from '@ui/media/icons/Plus';
import { Skeleton } from '@ui/feedback/Skeleton';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Percent03 } from '@ui/media/icons/Percent03';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { Currency, Opportunity, InternalStage } from '@graphql/types';
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
    const store = useStore();
    const [searchParams] = useSearchParams();
    const [isHover, setIsHover] = useState<boolean>();
    const { open, onOpen, onToggle } = useDisclosure();
    const viewDef = store.tableViewDefs.getById(
      store.tableViewDefs.opportunitiesPreset ?? '',
    );
    const column = viewDef?.value.columns.find((c) => c.columnId === columnId);

    const cards = store.opportunities.toComputedArray((arr) => {
      arr = arr
        .filter(
          (opp) =>
            opp.value.internalType === 'NBO' &&
            filterFns.every((fn) => fn(opp.value)),
        )
        .sort(
          (a, b) =>
            new Date(b.value.metadata.created).valueOf() -
            new Date(a.value.metadata.created).valueOf(),
        );

      if (searchParams.has('search')) {
        const search = searchParams.get('search')?.toLowerCase() ?? '';

        if (!search) return arr;

        arr = arr.filter((opp) => {
          return (
            opp.value.name.toLowerCase().includes(search) ||
            opp.organization?.value?.name.toLowerCase().includes(search) ||
            (
              opp.value.owner?.name ||
              [opp.value.owner?.firstName, opp.value.owner?.lastName].join(' ')
            )
              ?.toLowerCase()
              .includes(search)
          );
        });
      }

      return arr;
    });

    const currency = store.settings.tenant.value?.baseCurrency || Currency.Usd;
    const totalSum = formatCurrency(
      cards.reduce((acc, card) => acc + card.value.maxAmount, 0),
      0,
      currency,
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

    const handleCreateDraft = () => {
      store.ui.commandMenu.toggle('ChooseOpportunityOrganization', {
        ids: [],
        meta: { stage },
        entity: 'Opportunity',
      });
    };

    return (
      <div
        onMouseEnter={() => setIsHover(true)}
        onMouseLeave={() => setIsHover(false)}
        data-test={`kanban-column-${column?.name}`}
        className='flex flex-col flex-shrink-0 w-72 bg-gray-100 rounded h-full'
      >
        <div className='flex items-center justify-between p-3 pb-0 bg-gray-100 rounded-t-[4px]'>
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
                  onClick={(e) => (e.target as HTMLInputElement).select()}
                  className={cn(
                    'h-auto font-semibold min-h-[unset]',
                    !canEdit && 'cursor-not-allowed',
                  )}
                />
              </Tooltip>
              <Tooltip asChild label='Create new opportunity'>
                {canEdit && (
                  <IconButton
                    size='xxs'
                    icon={<Plus />}
                    onClick={handleCreateDraft}
                    aria-label='Add new opportunity'
                    dataTest={`add-opp-plus-${column?.name}`}
                    className={cn(
                      isHover ? 'opacity-100' : 'opacity-0',
                      canEdit && 'mr-2',
                    )}
                  />
                )}
              </Tooltip>
              {canEdit && (
                <Menu>
                  <MenuButton asChild>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      icon={<DotsVertical />}
                      aria-label='Column options'
                      dataTest={`opp-three-dots-menu-${column?.name}`}
                    />
                  </MenuButton>
                  <MenuList>
                    <MenuItem onClick={handleCreateDraft}>
                      <Plus />
                      Add new opportunity
                    </MenuItem>
                    <MenuItem onClick={onOpen}>
                      <Percent03 />
                      <span>Set win probability</span>
                    </MenuItem>
                  </MenuList>
                </Menu>
              )}
            </div>
            <span
              data-test={`card-sum-length-${column?.name}`}
              className={cn('w-full text-sm font-medium text-gray-500')}
            >
              {`${totalSum} • ${cards.length}`}
            </span>
          </div>
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
              data-test={`kanban-column-${column?.name}-cards`}
              className={cn('flex flex-col pb-2 p-3 overflow-auto', {
                'bg-gray-100': dropSnapshot?.isDraggingOver,
              })}
              {...dropProvided.droppableProps}
              style={{ height: 'calc(100% - 40px)' }}
            >
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

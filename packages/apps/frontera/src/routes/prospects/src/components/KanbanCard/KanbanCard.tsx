import { forwardRef } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Currency } from '@graphql/types';
import { User01 } from '@ui/media/icons/User01';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building06 } from '@ui/media/icons/Building06';
import { MaskedInput } from '@ui/form/Input/MaskedInput';
import { currencySymbol } from '@shared/util/currencyOptions';

interface DraggableKanbanCardProps {
  index: number;
  card: OpportunityStore;
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
  card: OpportunityStore;
  noPointerEvents?: boolean;
  provided?: DraggableProvided;
  snapshot?: DraggableStateSnapshot;
}

export const KanbanCard: React.FC<KanbanCardProps> = observer(
  ({ card, provided, snapshot, noPointerEvents }) => {
    const store = useStore();
    const navigate = useNavigate();

    if (!card.value.metadata.id) return null;

    const organization = card.organization;
    const logo = organization?.value.icon;

    const symbol =
      currencySymbol[store.settings.tenant.value?.baseCurrency ?? Currency.Usd];

    return (
      <div
        tabIndex={0}
        ref={provided?.innerRef}
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
        <div className='flex flex-col w-full items-start gap-2'>
          <div
            className='flex items-center gap-2'
            onMouseUp={() => {
              navigate(
                `/organization/${card.value?.organization?.metadata.id}/`,
              );
            }}
          >
            <Avatar
              name={`${card.value?.name}`}
              size='xs'
              icon={<Building06 className='text-primary-500 size-3' />}
              className='w-5 h-5 min-w-5'
              src={logo || undefined}
              variant='outlineSquare'
            />
            <nav className='text-sm font-medium shadow-none p-0 no-underline hover:no-underline focus:no-underline  line-clamp-1'>
              {card.value?.name}
            </nav>
          </div>

          <div className='flex items-center gap-2'>
            <Tooltip label={`${organization?.value?.owner?.name}`}>
              <Avatar
                name={`${organization?.value?.owner?.name}`}
                textSize='xs'
                size='xs'
                icon={<User01 className='text-primary-500 size-3' />}
                className={cn(
                  'w-5 h-5 min-w-5',
                  organization?.value?.owner?.profilePhotoUrl
                    ? ''
                    : 'border border-primary-200 text-xs',
                )}
                src={organization?.value?.owner?.profilePhotoUrl || ''}
              />
            </Tooltip>

            <MaskedInput
              variant='unstyled'
              size='xs'
              blocks={{
                num: {
                  mask: Number,
                  scale: 2,
                  thousandsSeparator: ',',
                  normalizeZeros: true,
                  padFractionalZeros: true,
                  radix: '.',
                },
              }}
              mask={`${symbol}num`}
              defaultValue={card.value.maxAmount.toString()}
              onAccept={(_, instance) => {
                card.update((value) => {
                  value.maxAmount = instance._unmaskedValue
                    ? parseFloat(instance._unmaskedValue)
                    : 0;

                  return value;
                });
              }}
            />
          </div>
        </div>
      </div>
    );
  },
);

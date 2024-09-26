import { useNavigate } from 'react-router-dom';
import { useRef, useState, forwardRef } from 'react';

import { useMergeRefs } from 'rooks';
import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Clock } from '@ui/media/icons/Clock';
import { useStore } from '@shared/hooks/useStore';
import { Divider } from '@ui/presentation/Divider';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { InternalStage } from '@shared/types/__generated__/graphql.types';

import { Owner } from './Owner';
import { MoreMenu } from './MoreMenu';
import { NextSteps } from './NextSteps';
import { ArrEstimate } from './ArrEstimate';
import { OpportunityName } from './OpportunityName';

interface DraggableKanbanCardProps {
  index: number;
  onBlur: () => void;
  isFocused?: boolean;
  card: OpportunityStore;
  noPointerEvents?: boolean;
  onFocus: (id: string) => void;
}

export const DraggableKanbanCard = forwardRef<
  HTMLDivElement,
  DraggableKanbanCardProps
>(({ card, index, noPointerEvents, onBlur, onFocus, isFocused }, _ref) => {
  return (
    <Draggable index={index} draggableId={card?.value?.metadata.id}>
      {(provided, snapshot) => {
        return (
          <KanbanCard
            card={card}
            onBlur={onBlur}
            onFocus={onFocus}
            provided={provided}
            snapshot={snapshot}
            isFocused={isFocused}
            noPointerEvents={noPointerEvents}
          />
        );
      }}
    </Draggable>
  );
});

interface KanbanCardProps {
  onBlur: () => void;
  isFocused?: boolean;
  card: OpportunityStore;
  noPointerEvents?: boolean;
  provided?: DraggableProvided;
  onFocus: (id: string) => void;
  snapshot?: DraggableStateSnapshot;
}

export const KanbanCard = observer(
  ({
    card,
    onBlur,
    onFocus,
    provided,
    snapshot,
    isFocused,
    noPointerEvents,
  }: KanbanCardProps) => {
    const store = useStore();
    const navigate = useNavigate();
    const containerRef = useRef<HTMLDivElement>(null);
    const nextStepsRef = useRef<HTMLTextAreaElement>(null);
    const mergedRef = useMergeRefs(provided?.innerRef, containerRef);
    const [showNextSteps, setShowNextSteps] = useState(!!card.value.nextSteps);

    const organization = card.organization;
    const logo = organization?.value.icon || organization?.value.logo;
    const daysInStage = card.value?.stageLastUpdated
      ? DateTimeUtils.differenceInDays(
          new Date().toISOString(),
          card.value?.stageLastUpdated,
        )
      : 0;

    const cardStage = card.value.internalStage;

    if (!card.value.metadata.id) return null;

    const handleNextStepsClick = () => {
      if (card.value.nextSteps) {
        setShowNextSteps(false);
        card.update((value) => {
          value.nextSteps = '';

          return value;
        });
      } else {
        setShowNextSteps(true);
        setTimeout(() => nextStepsRef.current?.focus(), 10);
      }
    };

    useOutsideClick({
      handler: () => {
        if (store.ui.commandMenu.isOpen) return;
        onBlur();
      },
      enabled: isFocused,
      ref: containerRef,
    });

    return (
      <div
        ref={mergedRef}
        {...provided?.draggableProps}
        {...provided?.dragHandleProps}
        onMouseEnter={() => onFocus(card.id)}
        className={cn(
          'group/kanbanCard  relative flex flex-col items-start px-3 pb-3 pt-[6px] mb-2 bg-white rounded-lg border border-gray-200 shadow-xs hover:shadow-lg focus:border-primary-500 transition-all duration-200 ease-in-out',
          {
            '!shadow-lg cursor-grabbing': snapshot?.isDragging,
            'pointer-events-none': noPointerEvents,
            'border-gray-400': isFocused,
          },
        )}
      >
        <div className='flex flex-col w-full items-start gap-2'>
          <div className='flex gap-2 w-full justify-between items-start'>
            <div className='flex gap-2 items-center'>
              <div className='flex flex-col'>
                <OpportunityName opportunityId={card.id} />
                <p
                  className='text-sm text-gray-500 p-0 hover:text-gray-700 hover:cursor-pointer'
                  onClick={() => {
                    logo &&
                      navigate(
                        `/organization/${card.value?.organization?.metadata.id}/`,
                      );
                  }}
                >
                  {organization?.value.name
                    ? organization.value.name
                    : 'Organization loading...'}
                </p>
              </div>
            </div>

            <MoreMenu
              hasNextSteps={!!card.value.nextSteps}
              onNextStepsClick={handleNextStepsClick}
            />
          </div>

          <div className='flex items-center gap-2 w-full'>
            <Owner opportunityId={card.id} ownerId={card.owner?.id} />

            <div className='flex items-center justify-between w-full mb-[-4px]'>
              <Tooltip label='ARR estimate'>
                <div>
                  <ArrEstimate opportunityId={card.id} />
                </div>
              </Tooltip>

              <Tooltip label='Time in stage'>
                {cardStage === InternalStage.Open && (
                  <div className='flex items-center'>
                    <Clock className='text-gray-500 size-3 mr-1' />
                    <span className='text-nowrap text-xs items-center'>
                      {`${daysInStage} ${daysInStage === 1 ? 'day' : 'days'}`}
                    </span>
                  </div>
                )}
              </Tooltip>
            </div>
          </div>
        </div>
        {(card.value.nextSteps || showNextSteps) && (
          <>
            <Divider className='mt-3 mb-2' />
            <NextSteps
              opportunityId={card.id}
              textareaRef={nextStepsRef}
              onToggle={setShowNextSteps}
            />
          </>
        )}
      </div>
    );
  },
);

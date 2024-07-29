import { useNavigate } from 'react-router-dom';
import { useRef, useState, forwardRef } from 'react';

import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Clock } from '@ui/media/icons/Clock';
import { User01 } from '@ui/media/icons/User01';
import { IconButton } from '@ui/form/IconButton';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Divider } from '@ui/presentation/Divider';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building06 } from '@ui/media/icons/Building06';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

import { NextSteps } from './NextSteps';
import { ArrEstimate } from './ArrEstimate';
import { OpportunityName } from './OpportunityName';

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
    const navigate = useNavigate();
    const nextStepsRef = useRef<HTMLTextAreaElement>(null);
    const [showNextSteps, setShowNextSteps] = useState(!!card.value.nextSteps);

    if (!card.value.metadata.id) return null;

    const organization = card.organization;
    const logo = organization?.value.icon;
    const daysInStage = card.value?.stageLastUpdated
      ? DateTimeUtils.differenceInDays(
          new Date().toISOString(),
          card.value?.stageLastUpdated,
        )
      : 0;

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

    return (
      <div
        ref={provided?.innerRef}
        {...provided?.draggableProps}
        {...provided?.dragHandleProps}
        className={cn(
          'group/kanbanCard !cursor-pointer relative flex flex-col items-start px-3 pb-2 pt-2 mb-2 bg-white rounded-lg border border-gray-200 shadow-xs hover:shadow-lg focus:border-primary-500 transition-all duration-200 ease-in-out',
          {
            '!shadow-lg cursor-grabbing': snapshot?.isDragging,
            'pointer-events-none': noPointerEvents,
          },
        )}
      >
        <div className='flex flex-col w-full items-start gap-2'>
          <div className='flex items-center gap-2 w-full justify-between'>
            <div className='flex gap-2 items-center'>
              <Avatar
                size='xs'
                variant='outlineSquare'
                src={logo || undefined}
                className='w-5 h-5 min-w-5'
                name={`${card.value?.name}`}
                icon={<Building06 className='text-primary-500 size-3' />}
                onMouseUp={() => {
                  navigate(
                    `/organization/${card.value?.organization?.metadata.id}/`,
                  );
                }}
              />
              <OpportunityName opportunityId={card.id} />
            </div>
            <Menu>
              <MenuButton asChild>
                <IconButton
                  size='xxs'
                  variant='ghost'
                  icon={<DotsVertical />}
                  aria-label='more options'
                />
              </MenuButton>

              <MenuList>
                <MenuItem onClick={handleNextStepsClick}>
                  {card.value.nextSteps ? <Trash01 /> : <ArrowsRight />}
                  {card.value.nextSteps ? 'Remove next step' : 'Add next step'}
                </MenuItem>
              </MenuList>
            </Menu>
          </div>

          <div className='flex items-center gap-2 w-full'>
            <Tooltip label={`${organization?.value?.owner?.name}`}>
              <Avatar
                size='xs'
                textSize='xs'
                name={`${organization?.value?.owner?.name}`}
                icon={<User01 className='text-primary-500 size-3' />}
                src={organization?.value?.owner?.profilePhotoUrl ?? ''}
                className={cn(
                  'w-5 h-5 min-w-5',
                  organization?.value?.owner?.profilePhotoUrl
                    ? ''
                    : 'border border-primary-200 text-xs',
                )}
              />
            </Tooltip>

            <div className='flex items-center justify-between w-full'>
              <ArrEstimate opportunityId={card.id} />

              <Clock className='text-gray-500 size-4 mr-1' />
              <span className='text-nowrap text-xs items-center'>
                {`${daysInStage} ${daysInStage === 1 ? 'day' : 'days'}`}
              </span>
            </div>
          </div>
        </div>
        {(card.value.nextSteps || showNextSteps) && (
          <>
            <Divider className='mt-1 mb-2' />
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

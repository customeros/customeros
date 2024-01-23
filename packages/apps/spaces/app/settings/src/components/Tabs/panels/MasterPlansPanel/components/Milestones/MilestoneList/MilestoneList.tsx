import { useRef } from 'react';

import { useDroppable } from '@dnd-kit/core';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';

import { Milestone } from './Milestone';
import { MilestoneDatum } from '../types';

interface MilestoneListProps {
  title: string;
  emptyText: string;
  droppeableId: string;
  milestones: MilestoneDatum[];
  onCreateMilestone: () => void;
  openMilestoneId: string | null;
  onToggleMilestone: (id: string) => void;
  onRemoveMilestone: (id: string) => void;
  onDuplicateMilestone: (id: string) => void;
  onMakeMilestoneOptional?: (id: string) => void;
  onSyncMilestone: (milestone: MilestoneDatum) => void;
}

export const MilestoneList = ({
  title,
  emptyText,
  milestones,
  droppeableId,
  openMilestoneId,
  onSyncMilestone,
  onCreateMilestone,
  onRemoveMilestone,
  onToggleMilestone,
  onDuplicateMilestone,
  onMakeMilestoneOptional,
}: MilestoneListProps) => {
  const shouldFocusNameRef = useRef(false);

  return (
    <>
      <Flex align='center' justify='space-between' mb='2'>
        <Text fontSize='sm' fontWeight='semibold'>
          {title}
        </Text>

        <IconButton
          size='xs'
          variant='ghost'
          onClick={() => {
            onCreateMilestone();
            shouldFocusNameRef.current = true;
          }}
          aria-label='Create milestone'
          icon={<Plus color='gray.400' />}
        />
      </Flex>
      <VStack>
        {milestones.length ? (
          milestones?.map((milestone, idx, arr) => (
            <Milestone
              key={milestone.id}
              milestone={milestone}
              onSync={onSyncMilestone}
              onToggle={onToggleMilestone}
              onRemove={onRemoveMilestone}
              isLast={idx === arr.length - 1}
              onDuplicate={onDuplicateMilestone}
              shouldFocusNameRef={shouldFocusNameRef}
              onMakeOptional={onMakeMilestoneOptional}
              isOpen={openMilestoneId === milestone.id}
            />
          ))
        ) : (
          <Empty emptyText={emptyText} droppableId={droppeableId} />
        )}
      </VStack>
    </>
  );
};

const Empty = ({
  emptyText,
  droppableId,
}: {
  emptyText: string;
  droppableId: string;
}) => {
  const { setNodeRef } = useDroppable({
    id: droppableId,
  });

  return (
    <Flex
      p='3'
      ref={setNodeRef}
      align='center'
      w='full'
      gap='2'
      bg='white'
      border='1px dashed'
      borderRadius='8px'
      borderColor='gray.300'
    >
      <Plus color='gray.400' />
      <Text color='gray.500'>{emptyText}</Text>
    </Flex>
  );
};

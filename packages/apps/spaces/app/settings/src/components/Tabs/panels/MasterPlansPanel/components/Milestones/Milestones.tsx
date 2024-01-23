import { useState } from 'react';

import { Active } from '@dnd-kit/core';

import { Flex } from '@ui/layout/Flex';
import { Divider } from '@ui/presentation/Divider';

import { MilestoneDatum } from './types';
import { MilestoneList } from './MilestoneList';
import { DndWrapper, DndSortableWrapper } from './DnD';
import { useMilestonesMethods } from './hooks/useMilestonesMethods';

interface MilestonesProps {
  milestones: MilestoneDatum[];
}

export const Milestones = ({ milestones }: MilestonesProps) => {
  const [openId, setOpenId] = useState<string | null>(null);
  const [active, setActive] = useState<Active | null>(null);

  const {
    allMilestones,
    onSyncMilestone,
    defaultMilestones,
    onCreateMilestone,
    onRemoveMilestone,
    optionalMilestones,
    onMilestonesChange,
    onDuplicateMilestone,
    onMakeMilestoneDefault,
    onMakeMilestoneOptional,
    onCreateOptionalMilestone,
  } = useMilestonesMethods({
    milestones,
  });

  const activeItem = allMilestones.find((item) => item.id === active?.id);

  const handleToggle = (id: string) => {
    setOpenId((prev) => (prev === id ? null : id));
  };

  return (
    <Flex overflowY='auto' flexDir='column' h={`calc(100vh - 64px)`}>
      <DndWrapper
        items={allMilestones}
        onActiveChange={setActive}
        onItemsChange={onMilestonesChange}
      >
        <DndSortableWrapper
          openId={openId}
          activeItem={activeItem}
          items={defaultMilestones}
        >
          <MilestoneList
            droppeableId='default'
            emptyText='Add default'
            openMilestoneId={openId}
            title='Default milestones'
            milestones={defaultMilestones}
            onToggleMilestone={handleToggle}
            onSyncMilestone={onSyncMilestone}
            onCreateMilestone={onCreateMilestone}
            onRemoveMilestone={onRemoveMilestone}
            onDuplicateMilestone={onDuplicateMilestone}
            onMakeMilestoneOptional={onMakeMilestoneOptional}
          />
        </DndSortableWrapper>

        <Divider mt='4' mb='2' borderColor='gray.300' variant='dashed' />

        <DndSortableWrapper
          openId={openId}
          activeItem={activeItem}
          items={optionalMilestones}
        >
          <MilestoneList
            droppeableId='optional'
            emptyText='Add optional'
            title='Optional milestones'
            openMilestoneId={openId}
            milestones={optionalMilestones}
            onToggleMilestone={handleToggle}
            onSyncMilestone={onSyncMilestone}
            onRemoveMilestone={onRemoveMilestone}
            onDuplicateMilestone={onDuplicateMilestone}
            onCreateMilestone={onCreateOptionalMilestone}
            onMakeMilestoneOptional={onMakeMilestoneDefault}
          />
        </DndSortableWrapper>
      </DndWrapper>
    </Flex>
  );
};

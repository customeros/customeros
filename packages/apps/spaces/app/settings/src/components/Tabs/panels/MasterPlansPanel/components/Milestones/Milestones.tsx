import { memo, useMemo, useState } from 'react';

import { produce } from 'immer';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import {
  Active,
  useSensor,
  DndContext,
  useSensors,
  DragOverlay,
  closestCenter,
  PointerSensor,
  DropAnimation,
  KeyboardSensor,
  defaultDropAnimationSideEffects,
} from '@dnd-kit/core';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';

import { Milestone } from './Milestone';
import { MilestoneDatum } from './types';

interface MilestonesProps {
  milestones: MilestoneDatum[];
}

export const Milestones = memo(({ milestones }: MilestonesProps) => {
  const [openId, setOpenId] = useState<string | null>(null);
  const [active, setActive] = useState<Active | null>(null);
  const [_milestones, setMilestones] = useState<MilestoneDatum[]>(milestones);

  const defaultMilestones = useMemo(() => {
    return _milestones.filter((m) => !m.optional);
  }, [_milestones]);
  const optionalMilestones = useMemo(() => {
    return _milestones.filter((m) => m.optional);
  }, [_milestones]);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  const activeItem = useMemo(
    () => _milestones.find((item) => item.id === active?.id),
    [active, _milestones],
  );

  const handleToggle = (id: string) => {
    setOpenId((prev) => (prev === id ? null : id));
  };

  const handleSyncMilestone = (milestone: MilestoneDatum) => {
    const nextMilestones = produce<MilestoneDatum[]>(_milestones, (draft) => {
      const idx = draft.findIndex((m) => m.id === milestone.id);
      if (idx === -1) return;

      draft[idx] = milestone;
    });

    setMilestones(nextMilestones);
  };

  const handleAddDefaultMilestone = () => {
    setMilestones((prev) => [
      ...prev,
      {
        id: `${prev?.length + 1}`,
        name: 'Unnamed milestone',
        items: [],
        durationHours: 24,
        optional: false,
        order: prev?.length + 1,
        retired: false,
      },
    ]);
  };
  const handleAddOptionalMilestone = () => {
    setMilestones((prev) => [
      ...prev,
      {
        id: `${prev?.length + 1}`,
        name: 'Unnamed milestone',
        items: [],
        durationHours: 24,
        optional: true,
        order: prev?.length + 1,
        retired: false,
      },
    ]);
  };

  const handleRemoveMilestone = (id: string) => {
    setMilestones((prev) => prev.filter((m) => m.id !== id));
  };

  const handleDuplicateMilestone = (id: string) => {
    const milestone = _milestones.find((m) => m.id === id);
    if (milestone) {
      setMilestones?.([
        ..._milestones,
        {
          ...milestone,
          id: `${_milestones?.length + 1}`,
        },
      ]);
    }
  };

  const handleMakeMilestoneOptional = (id: string) => {
    const nextMilestones = produce<MilestoneDatum[]>(_milestones, (draft) => {
      const idx = draft.findIndex((m) => m.id === id);
      if (idx === -1) return;

      draft[idx].optional = true;
    });

    setMilestones(nextMilestones);
  };

  return (
    <>
      <Flex align='center' justify='space-between' mb='2'>
        <Text fontSize='sm' fontWeight='semibold'>
          Default milestones
        </Text>

        <IconButton
          size='xs'
          variant='ghost'
          icon={<Plus color='gray.400' />}
          aria-label='Add Default Master Plan'
          onClick={handleAddDefaultMilestone}
        />
      </Flex>
      <VStack>
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={({ active }) => {
            setActive(active);
          }}
          onDragEnd={({ active, over }) => {
            if (over && active.id !== over?.id) {
              const activeIndex = _milestones.findIndex(
                ({ id }) => id === active.id,
              );
              const overIndex = _milestones.findIndex(
                ({ id }) => id === over.id,
              );

              setMilestones(arrayMove(_milestones, activeIndex, overIndex));
            }
            setActive(null);
          }}
          onDragCancel={() => {
            setActive(null);
          }}
        >
          <SortableContext
            items={_milestones}
            strategy={verticalListSortingStrategy}
          >
            {defaultMilestones?.map((milestone, idx, arr) => (
              <Milestone
                key={milestone.id}
                milestone={milestone}
                onToggle={handleToggle}
                onSync={handleSyncMilestone}
                isLast={idx === arr.length - 1}
                onRemove={handleRemoveMilestone}
                isOpen={openId === milestone.id}
                onDuplicate={handleDuplicateMilestone}
                onMakeOptional={handleMakeMilestoneOptional}
              />
            ))}
          </SortableContext>
          <SortableOverlay>
            {activeItem ? (
              <Milestone
                isActiveItem
                milestone={activeItem}
                isOpen={openId === activeItem.id}
              />
            ) : null}
          </SortableOverlay>
        </DndContext>
      </VStack>

      <Divider mt='4' mb='2' borderColor='gray.300' variant='dashed' />
      <Flex align='center' justify='space-between' mb='2'>
        <Text fontSize='sm' fontWeight='semibold'>
          Optional milestones
        </Text>

        <IconButton
          size='xs'
          variant='ghost'
          icon={<Plus color='gray.400' />}
          aria-label='Add Optional Master Plan'
          onClick={handleAddOptionalMilestone}
        />
      </Flex>

      <VStack>
        {optionalMilestones?.map((milestone, idx, arr) => (
          <Milestone
            key={milestone.id}
            milestone={milestone}
            onToggle={handleToggle}
            onSync={handleSyncMilestone}
            isLast={idx === arr.length - 1}
            onRemove={handleRemoveMilestone}
            isOpen={openId === milestone.id}
            onDuplicate={handleDuplicateMilestone}
          />
        ))}
      </VStack>
    </>
  );
});

const dropAnimationConfig: DropAnimation = {
  sideEffects: defaultDropAnimationSideEffects({
    styles: {
      active: {
        opacity: '0.4',
      },
    },
  }),
};

export function SortableOverlay({ children }: { children?: React.ReactNode }) {
  return (
    <DragOverlay dropAnimation={dropAnimationConfig}>{children}</DragOverlay>
  );
}

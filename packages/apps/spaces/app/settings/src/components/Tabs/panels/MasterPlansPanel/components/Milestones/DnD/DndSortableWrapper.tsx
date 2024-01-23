import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import {
  DragOverlay,
  DropAnimation,
  defaultDropAnimationSideEffects,
} from '@dnd-kit/core';

import { MilestoneDatum } from '../types';
import { Milestone } from '../MilestoneList';

interface DndSortableWrapperProps {
  openId?: string | null;
  items: MilestoneDatum[];
  children: React.ReactNode;
  activeItem?: MilestoneDatum;
}

export const DndSortableWrapper = ({
  items,
  openId,
  children,
  activeItem,
}: DndSortableWrapperProps) => {
  return (
    <>
      <SortableContext items={items} strategy={verticalListSortingStrategy}>
        {children}
      </SortableContext>
      <DndSortableOverlay>
        {activeItem && (
          <Milestone
            isActiveItem
            milestone={activeItem}
            isOpen={openId === activeItem.id}
          />
        )}
      </DndSortableOverlay>
    </>
  );
};

const dropAnimationConfig: DropAnimation = {
  duration: 0,
  sideEffects: defaultDropAnimationSideEffects({}),
};

export const DndSortableOverlay = ({
  children,
}: {
  children?: React.ReactNode;
}) => {
  return (
    <DragOverlay dropAnimation={dropAnimationConfig}>{children}</DragOverlay>
  );
};

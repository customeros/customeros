import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import {
  DragOverlay,
  DropAnimation,
  defaultDropAnimationSideEffects,
} from '@dnd-kit/core';

import { ColumnDef } from '@graphql/types';

interface DndSortableWrapperProps {
  items: ColumnDef[];
  children: React.ReactNode;
  renderActiveItem?: React.ReactNode;
}

export const DndSortableWrapper = ({
  items,
  children,
  renderActiveItem,
}: DndSortableWrapperProps) => {
  return (
    <>
      <SortableContext items={items} strategy={verticalListSortingStrategy}>
        {children}
      </SortableContext>
      <DndSortableOverlay>{renderActiveItem}</DndSortableOverlay>
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

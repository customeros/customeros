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
  DragEndEvent,
  closestCenter,
  PointerSensor,
  DragOverEvent,
  KeyboardSensor,
  DragStartEvent,
} from '@dnd-kit/core';

import { MilestoneDatum } from '../types';

interface DndWrapperProps {
  items: MilestoneDatum[];
  children: React.ReactNode;
  onActiveChange: (active: Active | null) => void;
  onItemsChange: (items: MilestoneDatum[]) => void;
}

export const DndWrapper = ({
  items,
  children,
  onItemsChange,
  onActiveChange,
}: DndWrapperProps) => {
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  const handleDragStart = ({ active }: DragStartEvent) =>
    onActiveChange(active);

  const handleDragEnd = ({ active, over }: DragEndEvent) => {
    if (over && active.id !== over?.id) {
      const activeIndex = items.findIndex(({ id }) => id === active.id);
      const overIndex = items.findIndex(({ id }) => id === over.id);

      const activeItem = items[activeIndex];
      const overItem = items[overIndex];
      if (activeItem?.optional && over.id === 'default') {
        const newItems = produce(items, (draft) => {
          draft[activeIndex].optional = false;
        });
        onItemsChange(arrayMove(newItems, activeIndex, overIndex));

        return;
      }

      if (!activeItem?.optional && over.id === 'optional') {
        const newItems = produce(items, (draft) => {
          draft[activeIndex].optional = true;
        });
        onItemsChange(arrayMove(newItems, activeIndex, overIndex));

        return;
      }

      if (activeItem?.optional !== overItem?.optional) {
        const newItems = produce(items, (draft) => {
          draft[activeIndex].optional = overItem.optional;
        });
        onItemsChange(arrayMove(newItems, activeIndex, overIndex));

        return;
      }

      onItemsChange(arrayMove(items, activeIndex, overIndex));
    }
    onActiveChange(null);
  };

  const handleDragOver = ({ active, over }: DragOverEvent) => {
    const activeContainerId = active?.data?.current?.sortable?.containerId;
    const overContainerId = over?.data?.current?.sortable?.containerId;

    const activeItem = items.find(({ id }) => id === active?.id);

    if (activeItem?.optional && over?.id === 'default') {
      const activeItemIndex = items.findIndex(({ id }) => id === active.id);
      if (activeItemIndex === -1) return;

      onItemsChange(
        produce(items, (draft) => {
          draft[activeItemIndex].optional = false;
        }),
      );

      return;
    }

    if (!activeItem?.optional && over?.id === 'optional') {
      const activeItemIndex = items.findIndex(({ id }) => id === active.id);
      if (activeItemIndex === -1) return;

      onItemsChange(
        produce(items, (draft) => {
          draft[activeItemIndex].optional = true;
        }),
      );

      return;
    }

    if (activeContainerId !== overContainerId) {
      const activeItemIndex = items.findIndex(({ id }) => id === active.id);
      const overItem = items.find(({ id }) => id === over?.id);

      if (activeItemIndex === -1) return;

      onItemsChange(
        produce(items, (draft) => {
          draft[activeItemIndex].optional = overItem?.optional ?? false;
        }),
      );
    }
  };

  const handleDragCancel = () => onActiveChange(null);

  return (
    <DndContext
      sensors={sensors}
      onDragEnd={handleDragEnd}
      onDragOver={handleDragOver}
      onDragStart={handleDragStart}
      onDragCancel={handleDragCancel}
      collisionDetection={closestCenter}
    >
      <SortableContext items={items} strategy={verticalListSortingStrategy}>
        {children}
      </SortableContext>
    </DndContext>
  );
};

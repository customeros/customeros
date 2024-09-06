import { forwardRef } from 'react';

import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { Pin02 } from '@ui/media/icons/Pin02';
import { ColumnViewType } from '@graphql/types';
import { MenuItem } from '@ui/overlay/Menu/Menu';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { HandleDrag } from '@ui/media/icons/HandleDrag';

interface DraggableColumnItemProps {
  index: number;
  label: string;
  visible: boolean;
  columnId?: number;
  helperText?: string;
  isDragDisabled?: boolean;
  noPointerEvents?: boolean;
  columnType: ColumnViewType;
  onCheck?: (columnId: number) => void;
}

export const DraggableColumnItem = forwardRef<
  HTMLDivElement,
  DraggableColumnItemProps
>(
  (
    {
      label,
      index,
      onCheck,
      visible,
      columnId,
      helperText,
      columnType,
      isDragDisabled,
      noPointerEvents,
    },
    _ref,
  ) => {
    return (
      <Draggable
        index={index}
        isDragDisabled={isDragDisabled}
        draggableId={String(columnId ?? Math.random())}
      >
        {(provided, snapshot) => {
          return (
            <ColumnItem
              label={label}
              onCheck={onCheck}
              visible={visible}
              columnId={columnId}
              provided={provided}
              snapshot={snapshot}
              columnType={columnType}
              helperText={helperText}
              isPinned={isDragDisabled}
              noPointerEvents={noPointerEvents}
            />
          );
        }}
      </Draggable>
    );
  },
);

interface ColumnItemProps {
  label: string;
  visible: boolean;
  columnId?: number;
  isPinned?: boolean;
  helperText?: string;
  noPointerEvents?: boolean;
  columnType: ColumnViewType;
  provided?: DraggableProvided;
  snapshot?: DraggableStateSnapshot;
  onCheck?: (columnId: number) => void;
}

export const ColumnItem = ({
  label,
  onCheck,
  visible,
  provided,
  snapshot,
  isPinned,
  columnId,
  columnType,
  helperText,
  noPointerEvents,
}: ColumnItemProps) => {
  return (
    <MenuItem
      ref={provided?.innerRef}
      onSelect={(e) => e.preventDefault()}
      className={cn(
        'group bg-white',
        snapshot?.isDragging && 'shadow-md',
        noPointerEvents && 'pointer-events-none',
      )}
      {...provided?.draggableProps}
      {...provided?.dragHandleProps}
    >
      <Checkbox
        size='sm'
        disabled={isPinned}
        isChecked={visible}
        onChange={() => {
          columnId && onCheck?.(columnId);
        }}
      />
      <div
        className={cn(
          'flex items-center w-full cursor-pointer',
          snapshot?.isDragging && 'cursor-grabbing',
        )}
      >
        <span
          data-test={`edit-col-${columnType}`}
          className={cn('flex-1', isPinned && 'text-gray-500')}
        >
          {label}
        </span>
        <span
          className={cn(
            'transition-opacity text-gray-500 select-none text-sm',
            isPinned ? 'opacity-100' : 'opacity-0 group-hover:opacity-100',
          )}
        >
          {isPinned ? 'Pinned' : helperText}
        </span>
        <div className='cursor-grab'>
          {isPinned ? (
            <Pin02 className='w-4 h-4 ml-2 text-gray-400' />
          ) : (
            <HandleDrag className='w-4 h-4 ml-2 text-gray-400' />
          )}
        </div>
      </div>
    </MenuItem>
  );
};

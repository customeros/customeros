'use client';
import { forwardRef } from 'react';

import {
  Draggable,
  DraggableProvided,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { cn } from '@ui/utils/cn';
import { ColumnViewType } from '@graphql/types';
import { MenuItem } from '@ui/overlay/Menu/Menu';
import { Checkbox } from '@ui/form/Checkbox/Checkbox2';
import { HandleDrag } from '@ui/media/icons/HandleDrag';

export const ColumnOption = forwardRef<
  HTMLDivElement,
  {
    index: number;
    label: string;
    visible: boolean;
    columnType: ColumnViewType;
    onCheck: (columnType: ColumnViewType) => void;
  }
>(({ columnType, label, index, onCheck, visible }, _ref) => {
  return (
    <Draggable
      index={index}
      draggableId={columnType}
      isDragDisabled={index === 0}
    >
      {(provided, snapshot) => {
        return (
          <ColumnOptionContent
            label={label}
            onCheck={onCheck}
            visible={visible}
            provided={provided}
            snapshot={snapshot}
            columnType={columnType}
          />
        );
      }}
    </Draggable>
  );
});

interface ColumnOptionContentProps {
  label: string;
  visible: boolean;
  columnType: ColumnViewType;
  provided: DraggableProvided;
  snapshot: DraggableStateSnapshot;
  onCheck: (columnType: ColumnViewType) => void;
}

export const ColumnOptionContent = ({
  label,
  provided,
  snapshot,
  onCheck,
  visible,
  columnType,
}: ColumnOptionContentProps) => {
  return (
    <MenuItem
      className={cn('group bg-white', snapshot.isDragging && 'shadow-md')}
      ref={provided.innerRef}
      onSelect={(e) => e.preventDefault()}
      {...provided.draggableProps}
    >
      <Checkbox
        className='mr-2'
        isChecked={visible}
        onChange={() => {
          onCheck(columnType);
        }}
      />
      <div className='flex items-center w-full'>
        <span className='flex-1'>{label}</span>
        <span className='opacity-0 group-hover:opacity-100 transition-opacity text-gray-500 select-none text-sm'>
          {'E.g. Zenith Contract'}
        </span>
        <div className='cursor-grab' {...provided.dragHandleProps}>
          <HandleDrag className='w-4 h-4 ml-2' />
        </div>
      </div>
    </MenuItem>
  );
};

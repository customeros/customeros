'use client';
import { forwardRef } from 'react';

import { Draggable } from '@hello-pangea/dnd';

import { ColumnDef } from '@graphql/types';
import { MenuItem } from '@ui/overlay/Menu/Menu';
import { Checkbox } from '@ui/form/Checkbox/Checkbox2';
import { HandleDrag } from '@ui/media/icons/HandleDrag';

export const ColumnOption = forwardRef<
  HTMLDivElement,
  {
    index: number;
    column: ColumnDef;
  }
>(({ column, index }, ref) => {
  return (
    <Draggable index={index} draggableId={column.id}>
      {(provided, snapshot) => {
        return (
          <MenuItem
            className='group'
            ref={provided.innerRef}
            onSelect={(e) => e.preventDefault()}
            {...provided.draggableProps}
            {...provided.dragHandleProps}
          >
            <Checkbox className='mr-2' />
            <div className='flex items-center w-full'>
              <span className='flex-1'>{column?.columnType?.name}</span>
              <span className='opacity-0 group-hover:opacity-100 transition-opacity text-gray-500 select-none'>
                {'eg: Sokeres'}
              </span>
              <div className='cursor-grab' {...provided.dragHandleProps}>
                <HandleDrag className='w-4 h-4 ml-2' />
              </div>
            </div>
          </MenuItem>
        );
      }}
    </Draggable>
  );
});

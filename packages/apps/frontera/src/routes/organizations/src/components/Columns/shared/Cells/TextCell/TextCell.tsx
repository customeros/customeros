import React, { useRef, ReactNode } from 'react';

import { TableCellTooltip } from '@ui/presentation/Table';

interface TextCellProps {
  text: string;
  leftIcon?: ReactNode;
}

export const TextCell = ({ text, leftIcon }: TextCellProps) => {
  const itemRef = useRef<HTMLDivElement>(null);

  if (!text) return <div className='text-gray-400'>Unknown</div>;

  return (
    <TableCellTooltip
      hasArrow
      label={text}
      align='start'
      side='bottom'
      targetRef={itemRef}
    >
      <div ref={itemRef} className='flex overflow-hidden'>
        {leftIcon && <div className='mr-1'>{leftIcon}</div>}
        <div className=' overflow-x-hidden overflow-ellipsis'>{text}</div>
      </div>
    </TableCellTooltip>
  );
};

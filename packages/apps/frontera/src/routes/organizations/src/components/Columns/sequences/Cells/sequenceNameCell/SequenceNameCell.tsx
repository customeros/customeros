import React, { useRef } from 'react';

import { TableCellTooltip } from '@ui/presentation/Table';

interface TextCellProps {
  text: string;
  unknownText?: string;
}

export const SequenceNameCell = ({
  text,
  unknownText = 'Unknown',
}: TextCellProps) => {
  const itemRef = useRef<HTMLDivElement>(null);

  if (!text) return <div className='text-gray-400'>{unknownText}</div>;

  return (
    <TableCellTooltip
      hasArrow
      label={text}
      align='start'
      side='bottom'
      targetRef={itemRef}
    >
      <div ref={itemRef} className='flex overflow-hidden'>
        <div className=' overflow-x-hidden overflow-ellipsis font-medium'>
          {text}
        </div>
      </div>
    </TableCellTooltip>
  );
};

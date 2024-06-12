import { useRef } from 'react';

import { TableCellTooltip } from '@ui/presentation/Table/TableCellTooltip';

export const IndustryCell = ({ value }: { value?: string }) => {
  const cellRef = useRef<HTMLDivElement>(null);

  if (!value) return <p className='text-gray-400'>Unknown</p>;

  return (
    <TableCellTooltip
      hasArrow
      align='start'
      side='bottom'
      label={value}
      targetRef={cellRef}
    >
      <p ref={cellRef} className='text-gray-700 cursor-default truncate'>
        {value}
      </p>
    </TableCellTooltip>
  );
};

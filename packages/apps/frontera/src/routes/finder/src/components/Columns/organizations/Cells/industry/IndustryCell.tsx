import { useRef } from 'react';

import { TableCellTooltip } from '@ui/presentation/Table/TableCellTooltip';

export const IndustryCell = ({
  value,
  enrichingStatus,
}: {
  value?: string;
  enrichingStatus: boolean;
}) => {
  const cellRef = useRef<HTMLDivElement>(null);

  if (!value)
    return (
      <p className='text-gray-400'>
        {enrichingStatus ? 'Enriching...' : 'Not set'}
      </p>
    );

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

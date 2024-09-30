import React from 'react';

import { observer } from 'mobx-react-lite';

interface FlowStatisticsCellProps {
  value?: number;
  total?: number;
}

export const FlowStatisticsCell = observer(
  ({ total, value }: FlowStatisticsCellProps) => {
    if (value === undefined || total === undefined) {
      return (
        <div className='text-gray-400' data-test={FlowStatisticsCell}>
          No data yet
        </div>
      );
    }

    // const percentage = total > 0 ? Math.round((value / total) * 100) : 0;

    return (
      <div className='flex items-center space-x-1 '>
        <span>{value}</span>
        {/*<span>({percentage}%)</span>*/}
      </div>
    );
  },
);

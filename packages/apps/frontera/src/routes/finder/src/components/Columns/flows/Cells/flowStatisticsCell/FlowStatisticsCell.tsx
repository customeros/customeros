import React from 'react';

import { observer } from 'mobx-react-lite';

interface FlowStatisticsCellProps {
  value?: number;
  total?: number;
  dataTest?: string;
}

export const FlowStatisticsCell = observer(
  ({ value, dataTest }: FlowStatisticsCellProps) => {
    if (typeof value !== 'number') {
      return (
        <div data-test={dataTest} className='text-gray-400'>
          No data yet
        </div>
      );
    }

    // const percentage = total > 0 ? Math.round((value / total) * 100) : 0;

    return (
      <div className='flex items-center space-x-1 '>
        <span data-test={`${dataTest}-in-flows-table`}>{value}</span>
        {/*<span>({percentage}%)</span>*/}
      </div>
    );
  },
);

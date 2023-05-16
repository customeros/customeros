import React from 'react';
import { Skeleton } from '@spaces/atoms/skeleton';

export const OrganizationTimelineSkeleton: React.FC = () => {
  const rows = Array(2)
    .fill('')
    .map((e, i) => i + 1);
  return (
    <div>
      {rows.map((row, i) => (
        <div
          key={`timeline-skeleton-${row}-${i}`}
          style={{ marginBottom: '8px' }}
        >
          <Skeleton height='60px' width='100%' />
        </div>
      ))}
    </div>
  );
};

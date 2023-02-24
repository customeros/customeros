import React from 'react';
import { Skeleton } from '../../../ui-kit/atoms/skeleton';

export const ListSkeleton = ({ id }: { id: string }) => {
  const rows = Array(4)
    .fill('')
    .map((e, i) => i + 1);
  return (
    <div>
      {rows.map((row) => (
        <div key={`${row}-${id}`}>
          <div>
            <Skeleton />
          </div>
          <div>
            <Skeleton />
          </div>
        </div>
      ))}
    </div>
  );
};

import { Position, Handle as FlowHandle } from '@xyflow/react';

import { cn } from '@ui/utils/cn.ts';

export const Handle = ({
  type,
  className,
}: {
  className?: string;
  type: 'target' | 'source';
}) => {
  return (
    <FlowHandle
      type={type}
      position={type === 'target' ? Position.Top : Position.Bottom}
      className={cn(`h-2 w-2 bg-white border-gray-700`, className, {
        'bg-transparent': type === 'target',
        'border-transparent': type === 'target',
      })}
    />
  );
};

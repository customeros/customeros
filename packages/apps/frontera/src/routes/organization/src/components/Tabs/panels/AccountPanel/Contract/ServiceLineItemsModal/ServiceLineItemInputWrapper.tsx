import { FC, ReactNode } from 'react';

import { cn } from '@ui/utils/cn';

export const ServiceLineItemInputWrapper: FC<{
  isDeleted: boolean;
  children: ReactNode;
  width: string | number;
}> = ({ children, isDeleted, width }) => {
  return (
    <div
      className={cn(
        isDeleted ? 'pointer-events-none' : 'pointer-events-auto',
        'text-sm',
      )}
      style={{ width }}
    >
      {children}
    </div>
  );
};

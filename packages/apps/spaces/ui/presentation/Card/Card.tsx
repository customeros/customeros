import * as React from 'react';

import { twMerge } from 'tailwind-merge';

import { useSlots } from '@ui/utils/hooks';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string;
  children: React.ReactNode;
}

export const CardHeader = ({ className, ...props }: CardProps) => {
  return (
    <div
      className={twMerge('flex flex-col space-y-1.5 p-6', className)}
      {...props}
    />
  );
};

export const CardContent = ({ children, className, ...props }: CardProps) => {
  return (
    <div className={twMerge('p-6 pt-0', className)} {...props}>
      {children}
    </div>
  );
};

export const CardFooter = ({ children, className, ...props }: CardProps) => {
  return (
    <div
      className={twMerge('flex items-center p-6 pt-0', className)}
      {...props}
    >
      {children}
    </div>
  );
};

export const Card = ({ children, className, ...props }: CardProps) => {
  const [contentSlot, footerSlot, headerSlot] = useSlots(
    children,
    CardContent,
    CardFooter,
    CardHeader,
  );

  return (
    <div
      className={twMerge(
        'rounded-xl border bg-card text-card-foreground shadow',
        className,
      )}
      {...props}
    >
      {headerSlot}
      {contentSlot}
      {footerSlot}
    </div>
  );
};

import * as React from 'react';

import { twMerge } from 'tailwind-merge';

import { useSlots } from '@ui/utils/hooks';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string;
  children: React.ReactNode;
}

export const CardHeader = ({ className, ...props }: CardProps) => {
  return <div className={twMerge('pb-1', className)} {...props} />;
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

export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ children, className, ...props }, ref) => {
    const [contentSlot, footerSlot, headerSlot] = useSlots(
      children,
      CardContent,
      CardFooter,
      CardHeader,
    );

    return (
      <div
        ref={ref}
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
  },
);

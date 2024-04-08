import React from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixCollapsible from '@radix-ui/react-collapsible';

interface CollapsibleProps extends RadixCollapsible.CollapsibleProps {
  className: string;
  children: React.ReactNode;
}

export const Collapsible = ({
  children,
  className,
  ...props
}: CollapsibleProps) => {
  return (
    <RadixCollapsible.Root className={twMerge('w-full', className)} {...props}>
      {children}
    </RadixCollapsible.Root>
  );
};

export const CollapsibleTrigger = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <RadixCollapsible.Trigger asChild={false}>
      {children}
    </RadixCollapsible.Trigger>
  );
};

export const CollapsibleContent = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <RadixCollapsible.Content className='data-[state="open"]:animate-collapseDown data-[state="closed"]:animate-collapseUp'>
      {children}
    </RadixCollapsible.Content>
  );
};

import React from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixCollapsible from '@radix-ui/react-collapsible';

interface CollapsibleProps extends RadixCollapsible.CollapsibleProps {
  className?: string;
  children: React.ReactNode;
}

export const CollapsibleRoot = ({
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

interface CollapsibleTriggerProps
  extends RadixCollapsible.CollapsibleTriggerProps {
  className?: string;
  children?: React.ReactNode;
}

export const CollapsibleTrigger = ({
  children,
  className,
  ...props
}: CollapsibleTriggerProps) => {
  return (
    <RadixCollapsible.Trigger className={twMerge(className)} {...props}>
      {children}
    </RadixCollapsible.Trigger>
  );
};

interface CollapsibleContentProps
  extends RadixCollapsible.CollapsibleContentProps {
  className?: string;
  children?: React.ReactNode;
}

export const CollapsibleContent = ({
  children,
  className,
  ...props
}: CollapsibleContentProps) => {
  return (
    <RadixCollapsible.Content
      {...props}
      className={twMerge(
        'data-[state="open"]:animate-collapseDown data-[state="closed"]:animate-collapseUp',
        className,
      )}
    >
      {children}
    </RadixCollapsible.Content>
  );
};

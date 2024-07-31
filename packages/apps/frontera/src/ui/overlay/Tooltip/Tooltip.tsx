import React from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixTooltip from '@radix-ui/react-tooltip';

export interface TooltipProps {
  open?: boolean;
  asChild?: boolean;
  tabIndex?: number;
  className?: string;
  hasArrow?: boolean;
  defaultOpen?: boolean;
  delayDuration?: number;
  label: React.ReactNode;
  children: React.ReactNode;
  align?: 'start' | 'end' | 'center';
  onOpenChange?: (open: boolean) => void;
  side?: 'top' | 'bottom' | 'left' | 'right';
}

export const Tooltip = ({
  side,
  open,
  label,
  align,
  children,
  hasArrow,
  tabIndex,
  className,
  defaultOpen,
  onOpenChange,
  delayDuration = 500,
  asChild = true,
}: TooltipProps) => {
  if (!label) return children;

  return (
    <RadixTooltip.Provider>
      <RadixTooltip.Root
        open={open}
        defaultOpen={defaultOpen}
        onOpenChange={onOpenChange}
        delayDuration={delayDuration}
      >
        <RadixTooltip.Trigger asChild={asChild} tabIndex={tabIndex}>
          {children}
        </RadixTooltip.Trigger>
        <RadixTooltip.Portal
          container={typeof window !== 'undefined' ? document?.body : null}
        >
          <RadixTooltip.Content
            side={side}
            align={align}
            sideOffset={5}
            className={twMerge(
              'z-[5000] data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade text-white select-none rounded-[4px] bg-gray-700 px-[15px] py-[10px] text-sm shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] will-change-[transform,opacity]',
              className,
            )}
          >
            {label}
            {hasArrow && <RadixTooltip.Arrow className='fill-gray-700' />}
          </RadixTooltip.Content>
        </RadixTooltip.Portal>
      </RadixTooltip.Root>
    </RadixTooltip.Provider>
  );
};

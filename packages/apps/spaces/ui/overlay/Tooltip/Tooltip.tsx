import React from 'react';

import { TooltipRootProps, Tooltip as ArkTooltip } from '@ark-ui/react';

import { cn } from '@ui/utils/cn';

type Placement =
  | 'top'
  | 'bottom'
  | 'left'
  | 'right'
  | 'bottom-start'
  | 'bottom-end'
  | 'top-start'
  | 'top-end'
  | 'left-start'
  | 'left-end'
  | 'right-start'
  | 'right-end';

interface TooltipProps extends TooltipRootProps {
  label: string;
  className?: string;
  hasArrow?: boolean;
  position: Placement;
  children: React.ReactNode;
}

export const Tooltip = ({
  children,
  position = 'bottom',
  label,
  hasArrow = false,
  className = '',
  ...props
}: TooltipProps) => {
  return (
    <ArkTooltip.Root
      closeDelay={150}
      openDelay={150}
      positioning={{ placement: position }}
      {...props}
    >
      {({ isOpen }) => (
        <>
          <ArkTooltip.Trigger>{children}</ArkTooltip.Trigger>
          <ArkTooltip.Positioner>
            <ArkTooltip.Content
              className={cn(
                'py-2 px-3 shadow-lg text-lg leading-7 rounded-md items-center overflow-hidden',
                {
                  'animate-in fade-in-15': isOpen,
                  'animate-out fade-out-15': !isOpen,
                },
                className,
              )}
            >
              {label}
            </ArkTooltip.Content>
            {hasArrow && (
              <>
                <ArkTooltip.Arrow
                  className={cn('data-[part=arrow]', {
                    'animate-in fade-in-55': isOpen,
                    'animate-out fade-out-20': !isOpen,
                  })}
                  style={
                    {
                      '--arrow-size': isOpen ? '8px' : '0',
                      '--arrow-background': isOpen
                        ? 'rgb(55 65 81)'
                        : 'transparent',
                    } as React.CSSProperties
                  }
                >
                  <ArkTooltip.ArrowTip />
                </ArkTooltip.Arrow>
              </>
            )}
          </ArkTooltip.Positioner>
        </>
      )}
    </ArkTooltip.Root>
  );
};

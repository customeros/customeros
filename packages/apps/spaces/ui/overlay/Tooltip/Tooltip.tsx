import { twMerge } from 'tailwind-merge';
import * as RadixTooltip from '@radix-ui/react-tooltip';

interface TooltipProps {
  label: string;
  className?: string;
  hasArrow?: boolean;
  children: React.ReactNode;
  align?: 'start' | 'end' | 'center';
  side?: 'top' | 'bottom' | 'left' | 'right';
}

export const Tooltip = ({
  label,
  children,
  className,
  align,
  side,
  hasArrow,
}: TooltipProps) => {
  return (
    <RadixTooltip.Provider>
      <RadixTooltip.Root>
        <RadixTooltip.Trigger asChild>{children}</RadixTooltip.Trigger>
        <RadixTooltip.Portal>
          <RadixTooltip.Content
            className={twMerge(
              'data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade text-white select-none rounded-[4px] bg-gray-700 px-[15px] py-[10px] text-[15px] leading-none shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] will-change-[transform,opacity]',
              className,
            )}
            side={side}
            align={align}
            sideOffset={5}
          >
            {label}
            {hasArrow && <RadixTooltip.Arrow className='fill-gray-700' />}
          </RadixTooltip.Content>
        </RadixTooltip.Portal>
      </RadixTooltip.Root>
    </RadixTooltip.Provider>
  );
};

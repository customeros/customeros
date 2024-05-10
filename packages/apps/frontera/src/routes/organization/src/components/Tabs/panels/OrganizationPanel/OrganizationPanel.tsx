import { useState, ReactNode, useEffect } from 'react';

import { cn } from '@ui/utils/cn';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

interface OrganizationPanelProps extends React.HTMLAttributes<HTMLDivElement> {
  title: string;
  bgImage?: string;
  withFade?: boolean;
  actionItem?: ReactNode;
  leftActionItem?: ReactNode;
  bottomActionItem?: ReactNode;
  shouldBlockPanelScroll?: boolean; // fix for https://linear.app/customer-os/issue/COS-619/scrollbar-overlaps-the-renewal-modals-in-safari
}
export const OrganizationPanel = ({
  bgImage,
  title,
  actionItem,
  leftActionItem,
  children,
  withFade = false,
  shouldBlockPanelScroll = false,
  bottomActionItem,
  ...props
}: OrganizationPanelProps) => {
  const [isMounted, setIsMounted] = useState(!withFade);

  useEffect(() => {
    if (!withFade) return;
    setIsMounted(true);
  }, []);

  return (
    <div
      className={cn('flex flex-1 flex-col h-full p-0 bg-no-repeat bg-contain')}
      style={{ backgroundImage: bgImage ? `url(${bgImage})` : '' }}
      {...props}
    >
      <div className='flex justify-between pt-4 pb-4 px-6'>
        <div className='flex items-center'>
          {leftActionItem && leftActionItem}
          <span className='text-lg text-gray-700 font-semibold'>{title}</span>
        </div>

        {actionItem && actionItem}
      </div>
      <ScrollAreaRoot>
        <ScrollAreaViewport>
          <div
            className={cn(
              isMounted ? 'opacity-100' : 'opacity-0',
              'flex flex-col space-y-2 justify-stretch w-full h-full px-6 pb-8 transition-opacity duration-300 ease-in-out',
            )}
          >
            {children}
          </div>
        </ScrollAreaViewport>
        <ScrollAreaScrollbar orientation='vertical'>
          <ScrollAreaThumb />
        </ScrollAreaScrollbar>
      </ScrollAreaRoot>
      {bottomActionItem && bottomActionItem}
    </div>
  );
};

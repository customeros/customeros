import { useState, ReactNode, useEffect } from 'react';

import { cn } from '@ui/utils/cn';

interface PanelContainerProps extends React.HTMLAttributes<HTMLDivElement> {
  title: string;
  bgImage?: string;
  withFade?: boolean;
  actionItem?: ReactNode;
}

export const PanelContainer = ({
  title,
  bgImage,
  children,
  actionItem,
  withFade = false,
  ...props
}: PanelContainerProps) => {
  const [isMounted, setIsMounted] = useState(!withFade);

  useEffect(() => {
    if (!withFade) return;
    setIsMounted(true);
  }, []);

  return (
    <div
      style={{ backgroundImage: bgImage ? `${bgImage}` : '' }}
      className='p-0 flex-1 flex flex-col h-full bg-no-repeat bg-contain'
      {...props}
    >
      <div className='flex justify-between pt-2 pb-4 px-6'>
        <p className='text-[16px] text-gray-700 font-semibold'>{title}</p>

        {actionItem}
      </div>

      <div
        className={cn(
          isMounted ? 'opacity-100' : 'opacity-0',
          'flex flex-col gap-2 w-full h-full px-6 pb-8 transition-opacity duration-300 ease-in-out',
        )}
      >
        {children}
      </div>
    </div>
  );
};

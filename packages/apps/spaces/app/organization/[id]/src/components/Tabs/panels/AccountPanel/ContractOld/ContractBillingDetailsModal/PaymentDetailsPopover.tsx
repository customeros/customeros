'use client';
import { useRouter } from 'next/navigation';
import React, { FC, ReactNode } from 'react';

import { cn } from '@ui/utils/cn';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface PaymentDetailsPopoverProps {
  content?: string;
  children: ReactNode;
  withNavigation?: boolean;
}

export const PaymentDetailsPopover: FC<PaymentDetailsPopoverProps> = ({
  withNavigation,
  content,
  children,
}) => {
  const { push } = useRouter();

  return (
    <Popover>
      <PopoverTrigger>
        <div className='w-full'>{children}</div>
      </PopoverTrigger>
      <PopoverContent
        className={cn(
          content ? 'block' : 'none',
          'w-fit bg-gray-700 text-white rounded-md text-sm border-none',
        )}
      >
        <div className='flex'>
          <p className='text-base mr-2'>{content}</p>

          {withNavigation && (
            <span
              className={'text-base underline text-white'}
              role='button'
              tabIndex={0}
              onClick={() => push('/settings?tab=billing')}
            >
              Go to Settings
            </span>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
};

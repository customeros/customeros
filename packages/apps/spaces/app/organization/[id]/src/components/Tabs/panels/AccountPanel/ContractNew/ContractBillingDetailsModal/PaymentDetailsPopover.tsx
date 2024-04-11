'use client';
import { useRouter } from 'next/navigation';
import React, { FC, ReactNode } from 'react';

import {
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

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
    <Popover placement='bottom' trigger='hover'>
      <PopoverTrigger>
        <div className='w-full'>{children}</div>
      </PopoverTrigger>
      <PopoverContent
        width='fit-content'
        bg='gray.700'
        color='white'
        borderRadius='md'
        boxShadow='none'
        fontSize='sm'
        border='none'
        display={content ? 'block' : 'none'}
      >
        <PopoverBody display='flex'>
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
        </PopoverBody>
      </PopoverContent>
    </Popover>
  );
};

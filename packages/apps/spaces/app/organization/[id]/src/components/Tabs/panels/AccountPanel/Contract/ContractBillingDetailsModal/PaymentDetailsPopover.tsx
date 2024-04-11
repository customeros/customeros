'use client';
import Link from 'next/link';
import React, { FC, ReactNode } from 'react';

import { Text } from '@ui/typography/Text';
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
          <Text mr={2} fontSize='md'>
            {content}
          </Text>

          {withNavigation && (
            <Text
              as={Link}
              fontSize='md'
              href='/settings?tab=billing'
              color='white'
              textDecoration='underline'
            >
              Go to Settings
            </Text>
          )}
        </PopoverBody>
      </PopoverContent>
    </Popover>
  );
};

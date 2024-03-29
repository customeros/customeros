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
      <PopoverTrigger>{children}</PopoverTrigger>
      <PopoverContent
        width='fit-content'
        bg='gray.700'
        color='white'
        mt={4}
        borderRadius='md'
        boxShadow='none'
        border='none'
        display={content ? 'block' : 'none'}
      >
        <PopoverBody display='flex'>
          <Text mr={2}>{content}</Text>

          {withNavigation && (
            <Text
              as={Link}
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

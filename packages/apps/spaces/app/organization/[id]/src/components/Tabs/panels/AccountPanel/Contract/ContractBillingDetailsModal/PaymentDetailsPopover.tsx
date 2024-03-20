'use client';
import Link from 'next/link';
import React, { FC, ReactNode } from 'react';

import { Text } from '@ui/typography/Text';
import {
  Popover,
  PopoverBody,
  PopoverArrow,
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
    <Popover placement='bottom-end' trigger='hover'>
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
        <PopoverArrow bg='gray.700' />

        <PopoverBody display='flex'>
          <Text mr={2}>{content}</Text>

          {withNavigation && (
            <Text as={Link} href='/settings?tab=billing' color='white'>
              Go to Settings
            </Text>
          )}
        </PopoverBody>
      </PopoverContent>
    </Popover>
  );
};

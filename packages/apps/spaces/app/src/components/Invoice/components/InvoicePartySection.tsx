'use client';

import React, { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

type InvoiceHeaderProps = {
  zip?: string;
  title: string;
  name?: string;
  email?: string;
  country?: string;
  locality?: string;
  isBlurred?: boolean;
  isFocused?: boolean;
  addressLine1?: string;
  addressLine2?: string;
};

export const InvoicePartySection: FC<InvoiceHeaderProps> = ({
  isBlurred,
  isFocused,
  zip = '',
  name = '',
  email = '',
  country = '',
  locality = '',
  addressLine1 = '',
  addressLine2 = '',
  title,
}) => (
  <Flex
    flexDir='column'
    flex={1}
    w={170}
    py={2}
    px={3}
    borderTop='1px solid'
    borderRight={title === 'From' ? 'none' : '1px solid'}
    borderBottom='1px solid'
    borderColor={'gray.300'}
    filter={isBlurred ? 'blur(2px)' : 'none'}
    transition='filter 0.25s ease-in-out'
    position='relative'
    sx={{
      '&:after': {
        content: '""',
        bg: 'transparent',
        border: '2px solid',
        position: 'absolute',
        top: 0,
        bottom: 0,
        left: 0,
        right: title === 'From' ? -4 : 0,
        opacity: isFocused ? 1 : 0,
        transition: 'opacity 0.25s ease-in-out',
      },
    }}
  >
    <Text fontWeight='semibold' mb={1} fontSize='sm'>
      {title}
    </Text>
    <Text fontSize='sm' fontWeight='medium' mb={1} lineHeight={1.2}>
      {name}
    </Text>

    <Text fontSize='sm' lineHeight={1.2} color='gray.500'>
      {addressLine1}
      <Text as='span' display='block' lineHeight={1.2}>
        {addressLine2}
      </Text>
    </Text>
    <Text fontSize='sm' lineHeight={1.2} color='gray.500'>
      {locality}
      {locality && zip && ', '} {zip}
    </Text>
    <Text fontSize='sm' lineHeight={1.2} color='gray.500'>
      {country}
    </Text>
    {email && (
      <Text fontSize='sm' lineHeight={1.2} color='gray.500'>
        {email}
      </Text>
    )}
  </Flex>
);

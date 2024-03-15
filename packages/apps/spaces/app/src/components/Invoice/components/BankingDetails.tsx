'use client';

import React, { FC } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { BankAccount } from '@graphql/types';

type InvoiceHeaderProps = {
  availableBankAccount?: Partial<BankAccount> | null;
};

export const BankingDetails: FC<InvoiceHeaderProps> = ({
  availableBankAccount,
}) => {
  return (
    <Flex flexDir='column' borderTop='1px solid' borderColor='gray.300' py={2}>
      <Text fontSize='xs' fontWeight='semibold'>
        Bank transfer
      </Text>
      <Flex justifyContent='space-between'>
        <Box>
          <Text fontSize='xs' fontWeight='medium'>
            Bank name
          </Text>
          <Text fontSize='xs' color='gray.500'>
            {availableBankAccount?.bankName || '-'}
          </Text>
        </Box>
        <Box>
          <Text fontSize='xs' fontWeight='medium'>
            Sort code
          </Text>
          <Text fontSize='xs' color='gray.500'>
            {availableBankAccount?.sortCode || '-'}
          </Text>
        </Box>
        <Box>
          <Text fontSize='xs' fontWeight='medium'>
            Account number
          </Text>
          <Text fontSize='xs' color='gray.500'>
            {availableBankAccount?.accountNumber || '-'}
          </Text>
        </Box>
        <Box w='25%'>
          <Text fontSize='xs' fontWeight='medium'>
            Other details
          </Text>
          <Text fontSize='xs' color='gray.500'>
            {availableBankAccount?.otherDetails || '-'}
          </Text>
        </Box>
      </Flex>
    </Flex>
  );
};

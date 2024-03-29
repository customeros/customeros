import React, { FC, useMemo } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Currency, BankAccount } from '@graphql/types';

type InvoiceHeaderProps = {
  currency?: string;
  availableBankAccount?: Partial<BankAccount> | null;
};

export const BankingDetails: FC<InvoiceHeaderProps> = ({
  availableBankAccount,
  currency,
}) => {
  const bankDetails: { label: string; value: string } = useMemo(() => {
    const details = {
      label: 'BIC/Swift',
      value: availableBankAccount?.bic || '-',
    };

    switch (currency) {
      case Currency.Gbp:
        details.label = 'Sort code';
        details.value = availableBankAccount?.sortCode || '-';
        break;
      case Currency.Usd:
        details.label = 'Routing Number';
        details.value = availableBankAccount?.routingNumber || '-';
        break;
      case Currency.Eur:
        details.label = 'BIC/Swift';
        details.value = availableBankAccount?.bic || '-';
        break;
      default:
        break;
    }

    return details;
  }, [currency, availableBankAccount]);

  const accountNumberLabel =
    currency === Currency.Eur ? 'IBAN' : 'Account number';
  const accountNumberValue =
    currency === Currency.Eur
      ? availableBankAccount?.iban
      : availableBankAccount?.accountNumber;

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
            {bankDetails.label}
          </Text>
          <Text fontSize='xs' color='gray.500'>
            {bankDetails.value}
          </Text>
        </Box>
        <Box>
          <Text fontSize='xs' fontWeight='medium'>
            {accountNumberLabel}
          </Text>
          <Text fontSize='xs' color='gray.500'>
            {accountNumberValue || '-'}
          </Text>
        </Box>
      </Flex>
    </Flex>
  );
};

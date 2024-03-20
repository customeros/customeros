'use client';

import React from 'react';
import Image from 'next/image';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { Invoice, BankAccount } from '@graphql/types';

import { ServicesTable } from './ServicesTable';
import logoCustomerOs from './assets/customeros-logo-tiny.png';
import {
  InvoiceHeader,
  InvoiceSummary,
  BankingDetails,
  InvoicePartySection,
} from './components';

// todo refactor, use generated type
type Address = {
  zip: string;
  email: string;
  name?: string;
  country?: string;
  locality: string;
  vatNumber?: string;
  addressLine1: string;
  addressLine2?: string;
};

type InvoiceProps = {
  tax: number;
  from: Address;
  total: number;
  dueDate: string;
  status?: string;
  subtotal: number;
  currency?: string;
  issueDate: string;
  billedTo: Address;
  amountDue?: number;
  note?: string | null;
  invoiceNumber: string;
  isBilledToFocused?: boolean;
  isInvoiceProviderFocused?: boolean;
  canPayWithBankTransfer?: boolean | null;
  lines: Partial<Invoice['invoiceLineItems']>;
  availableBankAccount?: Partial<BankAccount> | null;
};

export function Invoice({
  invoiceNumber,
  issueDate,
  dueDate,
  billedTo,
  from,
  lines,
  subtotal,
  tax,
  total,
  note,
  amountDue,
  status,
  isBilledToFocused,
  isInvoiceProviderFocused,
  currency = 'USD',
  canPayWithBankTransfer,
  availableBankAccount,
}: InvoiceProps) {
  const isInvoiceMetaSectionBlurred =
    isBilledToFocused || isInvoiceProviderFocused;

  return (
    <Flex
      px={4}
      flexDir='column'
      w='inherit'
      overflowY='auto'
      h='full'
      justifyContent='space-between'
      pb={4}
    >
      <Flex flexDir='column'>
        <Flex flexDir='column' mt={2}>
          <InvoiceHeader invoiceNumber={invoiceNumber} status={status} />

          <Flex
            mt={2}
            justifyContent='space-evenly'
            transition='filter 0.25s ease-in-out'
          >
            <Flex
              flexDir='column'
              flex={1}
              w={170}
              py={2}
              px={2}
              borderRight={'1px solid'}
              filter={isInvoiceMetaSectionBlurred ? 'blur(2px)' : 'none'}
              transition='filter 0.25s ease-in-out'
              borderTop='1px solid'
              borderBottom='1px solid'
              borderColor='gray.300'
            >
              <Text fontWeight='semibold' mb={1} fontSize='sm'>
                Issued
              </Text>
              <Text fontSize='sm' mb={4} color='gray.500'>
                {DateTimeUtils.format(
                  issueDate,
                  DateTimeUtils.dateWithAbreviatedMonth,
                )}
              </Text>
              <Text fontWeight='semibold' mb={1} fontSize='sm'>
                Due
              </Text>
              <Text fontSize='sm' color='gray.500'>
                {DateTimeUtils.format(
                  dueDate,
                  DateTimeUtils.dateWithAbreviatedMonth,
                )}
              </Text>
            </Flex>
            <InvoicePartySection
              title='Billed to'
              isBlurred={isInvoiceProviderFocused}
              isFocused={isBilledToFocused}
              zip={billedTo?.zip}
              name={billedTo?.name}
              email={billedTo?.email}
              country={billedTo?.country}
              locality={billedTo.locality}
              addressLine1={billedTo?.addressLine1}
              addressLine2={billedTo?.addressLine2}
              vatNumber={billedTo?.vatNumber}
            />
            <InvoicePartySection
              title='From'
              isBlurred={isBilledToFocused}
              isFocused={isInvoiceProviderFocused}
              zip={from?.zip}
              name={from?.name}
              email={from?.email}
              country={from?.country}
              locality={from?.locality}
              addressLine1={from?.addressLine1}
              addressLine2={from?.addressLine2}
              vatNumber={from?.vatNumber}
            />
          </Flex>
        </Flex>

        <Flex
          mt={4}
          flexDir='column'
          filter={isInvoiceMetaSectionBlurred ? 'blur(2px)' : 'none'}
          transition='filter 0.25s ease-in-out'
        >
          <ServicesTable services={lines} currency={currency} />
          <InvoiceSummary
            tax={tax}
            total={total}
            subtotal={subtotal}
            currency={currency}
            amountDue={amountDue}
            note={note}
          />
        </Flex>
      </Flex>

      <Box>
        {/*{canPayWithBankTransfer && availableBankAccount && (*/}
        <BankingDetails
          availableBankAccount={availableBankAccount}
          currency={currency}
        />
        {/*)}*/}
        <Flex
          alignItems='center'
          py={2}
          borderTop='1px solid'
          borderColor='gray.300'
        >
          <Box mr={2}>
            <Image
              src={logoCustomerOs}
              alt='CustomerOS'
              width={14}
              height={14}
            />
          </Box>
          <Text fontSize='xs' color='gray.500'>
            Powered by
            <Link
              color='gray.500'
              as='span'
              href='/'
              mx={1}
              textDecoration='underline'
            >
              CustomerOS
            </Link>
            - Revenue Intelligence for B2B hyperscalers
          </Text>
        </Flex>
      </Box>
    </Flex>
  );
}

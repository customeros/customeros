'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Invoice } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';

import { ServicesTable } from './ServicesTable';
import {
  InvoiceHeader,
  InvoiceSummary,
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
  lines: Partial<Invoice['invoiceLineItems']>;
  isDomesticBankingDetailsSectionFocused?: boolean;
  isInternationalBankingDetailsSectionFocused?: boolean;
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
  isInternationalBankingDetailsSectionFocused,
  isDomesticBankingDetailsSectionFocused,
}: InvoiceProps) {
  const isInvoiceMetaSectionBlurred =
    isBilledToFocused || isInvoiceProviderFocused;
  const bankingDetailsFocused =
    isDomesticBankingDetailsSectionFocused ||
    isInternationalBankingDetailsSectionFocused;
  const isServicesSectionBlurred =
    isInvoiceMetaSectionBlurred || bankingDetailsFocused;

  return (
    <Flex
      px={4}
      flexDir='column'
      w='inherit'
      overflowY='auto'
      h='full'
      justifyContent='space-between'
    >
      <Flex flexDir='column'>
        <Flex flexDir='column' mt={2}>
          <InvoiceHeader invoiceNumber={invoiceNumber} status={status} />

          <Flex
            mt={2}
            justifyContent='space-evenly'
            filter={bankingDetailsFocused ? 'blur(2px)' : 'none'}
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
          filter={isServicesSectionBlurred ? 'blur(2px)' : 'none'}
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

      {/*{(domesticBankingDetails || internationalBankingDetails) && (*/}
      {/*  <BankingDetails*/}
      {/*    isBlurred={Boolean(isInvoiceMetaSectionBlurred)}*/}
      {/*    domesticBankingDetails={domesticBankingDetails}*/}
      {/*    internationalBankingDetails={internationalBankingDetails}*/}
      {/*    isDomesticBankingDetailsSectionFocused={*/}
      {/*      isDomesticBankingDetailsSectionFocused*/}
      {/*    }*/}
      {/*    isInternationalBankingDetailsSectionFocused={*/}
      {/*      isInternationalBankingDetailsSectionFocused*/}
      {/*    }*/}
      {/*  />*/}
      {/*)}*/}
    </Flex>
  );
}

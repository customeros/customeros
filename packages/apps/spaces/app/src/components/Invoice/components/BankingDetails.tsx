'use client';

import React, { FC } from 'react';

import { Text } from '@ui/typography/Text';
import { Grid, GridItem } from '@ui/layout/Grid';

type InvoiceHeaderProps = {
  isBlurred: boolean;
  domesticBankingDetails?: string | null;
  internationalBankingDetails?: string | null;
  isDomesticBankingDetailsSectionFocused?: boolean;
  isInternationalBankingDetailsSectionFocused?: boolean;
};

export const BankingDetails: FC<InvoiceHeaderProps> = ({
  domesticBankingDetails,
  internationalBankingDetails,
  isBlurred,
  isDomesticBankingDetailsSectionFocused,
  isInternationalBankingDetailsSectionFocused,
}) => (
  <Grid
    templateColumns={
      !domesticBankingDetails || !internationalBankingDetails
        ? '1fr'
        : '1fr 1fr'
    }
    marginTop={6}
    minH={100}
    maxW={600}
    filter={isBlurred ? 'blur(2px)' : 'none'}
    transition='filter 0.25s ease-in-out'
  >
    {domesticBankingDetails && (
      <GridItem
        p={3}
        borderRight={internationalBankingDetails ? '1px solid' : 'none'}
        borderTop='1px solid'
        borderBottom='1px solid'
        borderColor='gray.300'
        filter={
          isInternationalBankingDetailsSectionFocused ? 'blur(2px)' : 'none'
        }
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
            right: 0,
            opacity: isDomesticBankingDetailsSectionFocused ? 1 : 0,
            transition: 'opacity 0.25s ease-in-out',
          },
        }}
      >
        <Text fontSize='xs' fontWeight='semibold'>
          Domestic Payments
        </Text>
        <Text fontSize='xs' whiteSpace='pre-wrap'>
          {domesticBankingDetails}
        </Text>
      </GridItem>
    )}

    {internationalBankingDetails && (
      <GridItem
        p={3}
        borderTop='1px solid'
        borderBottom='1px solid'
        borderColor='gray.300'
        filter={isDomesticBankingDetailsSectionFocused ? 'blur(2px)' : 'none'}
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
            right: domesticBankingDetails ? -4 : 0,
            opacity: isInternationalBankingDetailsSectionFocused ? 1 : 0,
            transition: 'opacity 0.25s ease-in-out',
          },
        }}
      >
        <Text fontSize='xs' fontWeight='semibold'>
          International Payments
        </Text>
        <Text fontSize='xs'>{internationalBankingDetails}</Text>
      </GridItem>
    )}
  </Grid>
);

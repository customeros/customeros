'use client';
import React from 'react';

import { useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody } from '@ui/presentation/Card';
import { Contract, RenewalSummary } from '@graphql/types';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getRenewalLikelihoodColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { useIsMutatingContract } from '@organization/src/components/Tabs/panels/AccountPanel/hooks/useIsMutatingContract';
import { useARRInfoModalContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';

interface ARRForecastProps {
  name: string;
  isInitialLoading?: boolean;
  contracts?: Contract[] | null;
  renewalSunnary?: RenewalSummary | null;
}

export const ARRForecast = ({
  isInitialLoading,
  renewalSunnary,
  name,
}: ARRForecastProps) => {
  const isRestoring = useIsRestoring();
  const { modal } = useARRInfoModalContext();

  const isUpdatingContract = useIsMutatingContract();
  const formattedMaxAmount = formatCurrency(
    renewalSunnary?.maxArrForecast ?? 0,
  );
  const formattedAmount = formatCurrency(renewalSunnary?.arrForecast ?? 0);

  const hasForecastChanged = formattedMaxAmount !== formattedAmount;

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        variant='ghost'
        bg='transparent'
        cursor='default'
        boxShadow='none'
        sx={{
          '& button': {
            opacity: 0,
            transition: 'opacity 0.2s linear',
          },
        }}
        _hover={{
          '& button': {
            opacity: 1,
          },
        }}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={getRenewalLikelihoodColor(
              renewalSunnary?.renewalLikelihood,
            )}
          >
            <CurrencyDollar />
          </FeaturedIcon>
          <Flex
            ml='5'
            w='full'
            align='center'
            columnGap={4}
            justify='space-between'
          >
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading
                  size='sm'
                  whiteSpace='nowrap'
                  fontWeight='semibold'
                  color='gray.700'
                  mr={2}
                >
                  ARR forecast
                </Heading>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Help'
                  onClick={(e) => {
                    e.stopPropagation();
                    modal.onOpen();
                  }}
                  icon={<Icons.HelpCircle color='gray.400' />}
                />
              </Flex>
            </Flex>

            <Flex flexDir='column'>
              <Heading
                fontSize='2xl'
                transition='opacity 0.25s ease-in'
                color={
                  isUpdatingContract && (!isInitialLoading || !isRestoring)
                    ? 'gray.400'
                    : 'gray.700'
                }
              >
                {formattedAmount}
              </Heading>
              {hasForecastChanged && !isUpdatingContract && (
                <Text
                  fontSize='sm'
                  textAlign='right'
                  textDecoration='line-through'
                >
                  {formattedMaxAmount}
                </Text>
              )}
            </Flex>
          </Flex>
        </CardBody>
      </Card>

      <InfoDialog
        isOpen={modal.isOpen}
        onClose={modal.onClose}
        onConfirm={modal.onClose}
        confirmButtonLabel='Got it'
        label='ARR forecast'
      >
        <Text fontSize='sm' fontWeight='normal' mb={4}>
          Annual Recurring Revenue (ARR) is the total amount of money you can
          expect to receive from
          <Text as='span' fontWeight='medium' mx={1}>
            {name ? name : `Unnamed`}
          </Text>
          for the next 12 months.
        </Text>
        <Text fontSize='sm' fontWeight='normal'>
          It includes all renewals but excludes one-time and per use services.
          Renewals are discounted based on the renewal likelihood
        </Text>
      </InfoDialog>
    </>
  );
};

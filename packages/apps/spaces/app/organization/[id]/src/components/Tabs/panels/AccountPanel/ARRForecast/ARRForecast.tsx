'use client';
import React from 'react';
import { useParams } from 'next/navigation';

import { useIsFetching, useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody } from '@ui/presentation/Card';
import { Contract, RenewalForecast } from '@graphql/types';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { useARRInfoModalContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';

interface ARRForecastProps {
  name: string;
  isInitialLoading?: boolean;
  contracts?: Contract[] | null;
  forecast?: RenewalForecast | null;
}

export const ARRForecast = ({
  isInitialLoading,
  forecast,
  contracts,
  name,
}: ARRForecastProps) => {
  const isRestoring = useIsRestoring();
  const id = useParams()?.id as string;
  const { modal } = useARRInfoModalContext();

  const queryKey = useGetContractsQuery.getKey({ id });

  const isFetching = useIsFetching({
    queryKey: queryKey,
  });
  const formattedMaxAmount = formatCurrency(forecast?.maxArr ?? 0);
  const formattedAmount = formatCurrency(forecast?.arr ?? 0);

  const hasForecastChanged = formattedMaxAmount !== formattedAmount;
  const iconColor = (() => {
    if (!contracts?.length) {
      return 'gray';
    }

    const array = contracts.flatMap((contract) =>
      contract.opportunities
        ? contract.opportunities.map((opportunity) =>
            opportunity.renewalLikelihood.toLowerCase(),
          )
        : [],
    );

    const likelihood = array.every((val) => val === array[0])
      ? array[0]
      : 'medium';

    const colorMap: Record<string, string> = {
      high: 'success',
      medium: 'warning',
      low: 'error',
      zero: 'zero',
    };

    return colorMap[likelihood] || 'gray';
  })();

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
          <FeaturedIcon size='md' minW='10' colorScheme={iconColor}>
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
              <Heading fontSize='2xl' color='gray.700'>
                {isFetching && (!isInitialLoading || !isRestoring)
                  ? 'Calculating...'
                  : formattedAmount}
              </Heading>
              {hasForecastChanged && !isFetching && (
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
            {name}
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

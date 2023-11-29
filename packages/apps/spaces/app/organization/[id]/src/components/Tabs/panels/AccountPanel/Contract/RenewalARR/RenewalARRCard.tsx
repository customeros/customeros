import React, { useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { Card, CardHeader } from '@ui/presentation/Card';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import {
  Opportunity,
  ContractRenewalCycle,
  RenewalLikelihoodProbability,
} from '@graphql/types';
import { RenewalDetailsModal } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/RenewalARR/RenewalDetailsModal';
import { useUpdateRenewalDetailsContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';

import { getARRColor } from '../../utils';

interface RenewalARRCardProps {
  hasEnded: boolean;
  startedAt: string;
  opportunity: Opportunity;
  renewCycle: ContractRenewalCycle;
}
export const RenewalARRCard = ({
  startedAt,
  hasEnded,
  renewCycle,
  opportunity,
}: RenewalARRCardProps) => {
  const { modal } = useUpdateRenewalDetailsContext();
  const [isLocalOpen, setIsLocalOpen] = useState(false);

  const differenceInMonths = DateTimeUtils.differenceInMonths(
    new Date().toISOString(),
    startedAt,
  );

  const hasRenewed =
    renewCycle === ContractRenewalCycle.AnnualRenewal
      ? differenceInMonths > 12
      : differenceInMonths > 1;

  const formattedMaxAmount = formatCurrency(opportunity.maxAmount ?? 0);
  const formattedAmount = formatCurrency(hasEnded ? 0 : opportunity.amount);

  const hasRewenewChanged = formattedMaxAmount !== formattedAmount;

  return (
    <>
      <Card
        px='4'
        py='3'
        w='full'
        my={2}
        size='lg'
        variant='outline'
        cursor='pointer'
        border='1px solid'
        borderColor='gray.200'
        position='relative'
        onClick={() => {
          modal.onOpen();
          setIsLocalOpen(true);
        }}
        width={hasRenewed ? 'calc(100% - .5rem)' : 'auto'}
        sx={
          hasRenewed
            ? {
                right: -2,
                '&:after': {
                  content: "''",
                  width: 2,
                  height: '80%',
                  left: '-9px',
                  top: '6px',
                  bg: 'white',
                  position: 'absolute',
                  borderTopLeftRadius: 'md',
                  borderBottomLeftRadius: 'md',
                  border: '1px solid',
                  borderColor: 'gray.200',
                },
              }
            : {}
        }
      >
        <CardHeader
          as={Flex}
          p='0'
          w='full'
          alignItems='center'
          justifyContent='center'
          gap={4}
        >
          <FeaturedIcon size='md' minW='10' colorScheme='primary'>
            <ClockFastForward />
          </FeaturedIcon>
          <Flex alignItems='center' justifyContent='space-between' w='full'>
            <Flex flexDir='column' gap='1px'>
              <Flex flex={1} alignItems='center'>
                <Heading size='sm' color='gray.700' noOfLines={1}>
                  Renewal ARR
                </Heading>

                {!hasEnded && opportunity.renewedAt && startedAt && (
                  <Text color='gray.500' ml={1} fontSize='sm'>
                    {DateTimeUtils.isToday(opportunity.renewedAt)
                      ? 'today'
                      : DateTimeUtils.timeAgo(opportunity.renewedAt, {
                          addSuffix: true,
                        })}
                  </Text>
                )}
              </Flex>

              {startedAt && (
                <Text w='full' color='gray.500' fontSize='sm' lineHeight={1}>
                  {!hasEnded ? (
                    <>
                      Likelihood{' '}
                      <Text
                        as='span'
                        fontWeight='medium'
                        color={`${getARRColor(
                          opportunity.renewalLikelihood as RenewalLikelihoodProbability,
                        )}.500`}
                        textTransform='capitalize'
                      >
                        {opportunity?.renewalLikelihood.toLowerCase()}
                      </Text>
                    </>
                  ) : (
                    'Closed lost'
                  )}
                </Text>
              )}
            </Flex>

            <Flex flexDir='column'>
              <Text fontWeight='semibold'>{formattedAmount}</Text>

              {hasRewenewChanged && (
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
        </CardHeader>
      </Card>
      <RenewalDetailsModal
        isOpen={modal.isOpen && isLocalOpen}
        onClose={() => {
          modal.onClose();
          setIsLocalOpen(false);
        }}
        data={opportunity}
      />
    </>
  );
};

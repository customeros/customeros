import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { ContractRenewalCycle } from '@graphql/types';
import { Card, CardHeader } from '@ui/presentation/Card';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { calculateNextRenewalDate } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
// import { RenewalDetailsModal } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/RenewalDetailsModal';
// import { useARRInfoModalContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';

interface RenewalARRCardProps {
  hasEnded: boolean;
  startedAt: string;
  renewCycle: ContractRenewalCycle;
}
export const RenewalARRCard = ({
  startedAt,
  hasEnded,
  renewCycle,
}: RenewalARRCardProps) => {
  const nextRenewal = calculateNextRenewalDate(startedAt, renewCycle);
  const differenceInMonths = DateTimeUtils.differenceInMonths(
    new Date().toISOString(),
    startedAt,
  );
  const hasRenewed =
    renewCycle === ContractRenewalCycle.AnnualRenewal
      ? differenceInMonths > 12
      : differenceInMonths > 1;

  return (
    <>
      <Card
        px='4'
        py='3'
        w='full'
        my={2}
        size='lg'
        variant='outline'
        cursor='default'
        border='1px solid'
        borderColor='gray.200'
        position='relative'
        // onClick={modal.onOpen}
        width={hasRenewed ? 'calc(100% - .5rem)' : 'auto'}
        sx={
          hasRenewed
            ? {
                right: -2,
                '&:after': {
                  content: "''",
                  width: 2,
                  height: '80%',
                  left: -2,
                  top: '7px',
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
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={
              'gray'
              // renewalForecast
              //   ? getFeatureIconColor(
              //       renewalProbability as
              //         | RenewalLikelihoodProbability
              //         | undefined,
              //     )
              //   :
            }
          >
            <ClockFastForward />
          </FeaturedIcon>
          <Flex
            alignItems='center'
            justifyContent='space-between'
            w='full'
            // mt={-4}
          >
            <Flex flex={1} alignItems='center'>
              <Heading size='sm' color='gray.700' noOfLines={1}>
                Renewal ARR
              </Heading>

              {!hasEnded && (
                <Text color='gray.500' ml={1} fontSize='sm'>
                  {DateTimeUtils.isToday(nextRenewal)
                    ? 'today'
                    : DateTimeUtils.timeAgo(nextRenewal, { addSuffix: true })}
                </Text>
              )}
            </Flex>
            {/*TODO swap with real data after integrating with BE*/}
            {/*<Text fontWeight='semibold'>$24,000</Text>*/}
          </Flex>
        </CardHeader>

        {/*<CardBody*/}
        {/*  as={Text}*/}
        {/*  w='full'*/}
        {/*  color='gray.500'*/}
        {/*  fontSize='sm'*/}
        {/*  p={0}*/}
        {/*  pl={14}*/}
        {/*  mt={-5}*/}
        {/*>*/}
        {/*  Likelihood{' '}*/}
        {/*  <Text as='span' color='success.500'>*/}
        {/*    High*/}
        {/*  </Text>*/}
        {/*</CardBody>*/}
      </Card>
      {/*<RenewalDetailsModal*/}
      {/*  isOpen={modal.isOpen}*/}
      {/*  onClose={modal.onClose}*/}
      {/*  data={{}}*/}
      {/*/>*/}
    </>
  );
};

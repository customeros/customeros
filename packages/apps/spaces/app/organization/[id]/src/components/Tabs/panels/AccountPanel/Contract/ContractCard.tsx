import { useForm } from 'react-inverted-form';
import React, { useRef, useState, useEffect } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { FormInput } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Edit03 } from '@ui/media/icons/Edit03';
import { UseDisclosureReturn } from '@ui/utils';
import { File02 } from '@ui/media/icons/File02';
import { FormSelect } from '@ui/form/SyncSelect';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { ChevronUp } from '@ui/media/icons/ChevronUp';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import { Contract, ContractRenewalCycle } from '@graphql/types';
import { Card, CardBody, CardFooter, CardHeader } from '@ui/presentation/Card';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { ServiceModal } from '@organization/src/components/Tabs/panels/AccountPanel/Services/ServiceModal';
// import { RenewalARRCard } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/RenewalARRCard';
import {
  ContractDTO,
  TimeToRenewalForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Contract.dto';
import { ContractStatusSelect } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/contractStatuses/ContractStatusSelect';

interface ContractCardProps {
  // todo use generated type after gql schema for service item is merged
  data: Contract;
  serviceModal: UseDisclosureReturn;
  //   & {
  // services: Array<{
  //   id: string;
  //   price: number;
  //   billed: string;
  //   quantity: number;
  // }>;
}

function getLabelFromValue(value: string): string | undefined {
  if (ContractRenewalCycle.AnnualRenewal === value) {
    return 'annually';
  }
  if (ContractRenewalCycle.MonthlyRenewal === value) {
    return 'monthly';
  }
}
export const ContractCard = ({ data, serviceModal }: ContractCardProps) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isExpanded, setIsExpanded] = useState(!data?.signedAt);
  const formId = 'contractForm';

  const defaultValues = ContractDTO.toForm(data);
  const { setDefaultValues } = useForm<TimeToRenewalForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [
    defaultValues.signedAt?.toISOString(),
    defaultValues.renewalCycle,
    defaultValues.endedAt?.toISOString(),
    defaultValues.serviceStartedAt?.toISOString(),
  ]);

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  function calculateNextRenewalDate(
    date: string,
    renewalCycle: ContractRenewalCycle,
    period: number,
  ): string | number {
    switch (renewalCycle) {
      case ContractRenewalCycle.AnnualRenewal:
        return DateTimeUtils.addYears(date, period).toISOString();
      default:
        return DateTimeUtils.addMonth(date, period).toISOString();
    }
  }

  return (
    <Card
      px='4'
      py='3'
      w='full'
      size='lg'
      variant='outline'
      cursor='default'
      border='1px solid'
      borderColor='gray.200'
      bg='gray.50'
      transition='all 0.2s ease-out'
    >
      <CardHeader
        as={Flex}
        p='0'
        pb={!data?.signedAt ? 2 : 0}
        w='full'
        flexDir='column'
      >
        <Flex justifyContent='space-between' w='full' flex={1}>
          <Heading size='sm' color='gray.700' noOfLines={1} w='fit-content'>
            {data?.name}

            {!data?.name && (
              <FormInput
                pointerEvents={isExpanded ? 'all' : 'none'}
                fontWeight='semibold'
                fontSize='inherit'
                name='name'
                formId={formId}
                placeholder='Contract name'
                _hover={{
                  borderBottom: !isExpanded && 'none',
                }}
              />
            )}
          </Heading>
          <Flex alignItems='center' gap={2} ml={4}>
            {data?.contractUrl && (
              <Tooltip label='Open contract url'>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Open contract'
                  icon={<File02 color='gray.400' />}
                  onClick={() =>
                    window.open(
                      getExternalUrl(
                        getExternalUrl(data.contractUrl as string),
                      ),
                      '_blank',
                      'noopener',
                    )
                  }
                />
              </Tooltip>
            )}

            <ContractStatusSelect />

            {(isExpanded || (!isExpanded && !data?.name)) && (
              <IconButton
                size='xs'
                variant='ghost'
                aria-label='Collapse'
                icon={isExpanded ? <ChevronUp /> : <ChevronDown />}
                onClick={() => setIsExpanded(!isExpanded)}
              />
            )}
          </Flex>
        </Flex>

        {!isExpanded && (
          <Button
            bg='transparent'
            _hover={{
              bg: 'transparent',
              svg: { opacity: 1, transition: 'opacity 0.2s linear' },
            }}
            sx={{ svg: { opacity: 0, transition: 'opacity 0.2s linear' } }}
            size='xs'
            fontSize='sm'
            fontWeight='normal'
            color='gray.500'
            p={0}
            alignItems='flex-start'
            justifyContent='flex-start'
            onClick={() => setIsExpanded(true)}
          >
            <Flex
              flexDir='column'
              alignItems='flex-start'
              justifyContent='center'
            >
              {!data?.signedAt && (
                <Text mt={-2}>No start date or services yet</Text>
              )}

              {data?.signedAt && (
                <>
                  <Text>
                    {data?.serviceStartedAt &&
                      DateTimeUtils.isFuture(data.serviceStartedAt) &&
                      `Service starts on ${DateTimeUtils.format(
                        data.serviceStartedAt,
                        DateTimeUtils.dateWithAbreviatedMonth,
                      )}`}
                  </Text>
                  <Text>
                    {data?.renewalCycle &&
                      !DateTimeUtils.isFuture(data.serviceStartedAt) &&
                      `Service renews ${getLabelFromValue(
                        data.renewalCycle,
                      )} on ${DateTimeUtils.format(
                        calculateNextRenewalDate(
                          data.serviceStartedAt,
                          data.renewalCycle,
                          1,
                        ),
                        DateTimeUtils.dateWithAbreviatedMonth,
                      )}`}
                  </Text>

                  {!data?.renewalCycle &&
                    data?.endedAt &&
                    DateTimeUtils.isFuture(data.endedAt) && (
                      <Text>
                        Ends on{' '}
                        {DateTimeUtils.format(
                          data.endedAt,
                          DateTimeUtils.dateWithAbreviatedMonth,
                        )}
                      </Text>
                    )}

                  {data?.renewalCycle &&
                    data?.endedAt &&
                    DateTimeUtils.isFuture(data.endedAt) && (
                      <Text>
                        Ends on{' '}
                        {DateTimeUtils.format(
                          data.endedAt,
                          DateTimeUtils.dateWithAbreviatedMonth,
                        )}
                      </Text>
                    )}

                  {data?.endedAt && !DateTimeUtils.isFuture(data.endedAt) && (
                    <Text>
                      Ended on{' '}
                      {DateTimeUtils.format(
                        data.endedAt,
                        DateTimeUtils.dateWithAbreviatedMonth,
                      )}
                    </Text>
                  )}
                </>
              )}
            </Flex>

            <Edit03 ml={1} color='gray.400' boxSize='3' />
          </Button>
        )}
      </CardHeader>
      {isExpanded && (
        <CardBody as={Flex} p='0' flexDir='column' w='full'>
          <Flex gap='4' mb={2}>
            <DatePicker
              label='Contract signed'
              placeholder='Signed date'
              formId={formId}
              name='signedAt'
              calendarIconHidden
              inset='120% auto auto 0px'
            />
            <DatePicker
              label='Contract ends'
              placeholder='End date'
              formId={formId}
              name='endedAt'
              calendarIconHidden
            />
          </Flex>
          <Flex gap='4'>
            <DatePicker
              label='Service starts'
              placeholder='Start date'
              formId={formId}
              name='serviceStartedAt'
              calendarIconHidden
              inset='120% auto auto 0px'
            />

            <FormSelect
              label='Contract renews'
              placeholder='Contract renews'
              isLabelVisible
              name='renewalCycle'
              formId={formId}
              options={billingFrequencyOptions}
            />
          </Flex>
        </CardBody>
      )}

      <CardFooter p='0' w='full' flexDir='column'>
        {/*{!!data?.services && (*/}
        {/*  <RenewalARRCard*/}
        {/*  // withMultipleServices={!!data?.services && data?.services?.length > 1}*/}
        {/*  />*/}
        {/*)}*/}

        <Flex w='full' alignItems='center' justifyContent='space-between'>
          <Text fontWeight='semibold' fontSize='sm'>
            No services
            {/*{!data?.services?.length ? 'No services' : 'Services'}*/}
          </Text>

          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Add service'
            color='gray.400'
            onClick={() => serviceModal.onOpen()}
            icon={<Plus boxSize='4' />}
          />
        </Flex>

        {/*{data?.services?.length && <ServicesList data={data?.services} />}*/}
      </CardFooter>
      <ServiceModal
        contractId={data.id}
        isOpen={serviceModal.isOpen}
        onClose={serviceModal.onClose}
      />
    </Card>
  );
};

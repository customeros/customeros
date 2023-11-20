import { useForm } from 'react-inverted-form';
import React, { useRef, useState, useEffect } from 'react';

import { debounce } from 'lodash';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { FormInput } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Edit03 } from '@ui/media/icons/Edit03';
import { FormSelect } from '@ui/form/SyncSelect';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { ChevronUp } from '@ui/media/icons/ChevronUp';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardBody, CardFooter, CardHeader } from '@ui/presentation/Card';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { useUpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  Contract,
  ContractStatus,
  ContractUpdateInput,
  ContractRenewalCycle,
} from '@graphql/types';

import { UrlInput } from './UrlInput';
import { Services } from './Services/Services';
import { RenewalARRCard } from './RenewalARRCard';
import { ContractDTO, TimeToRenewalForm } from './Contract.dto';
import { billingFrequencyOptions, calculateNextRenewalDate } from '../utils';
import { ContractStatusSelect } from './contractStatuses/ContractStatusSelect';

interface ContractCardProps {
  data: Contract;
  organizationId: string;
  organizationName: string;
}

function getLabelFromValue(value: string): string | undefined {
  if (ContractRenewalCycle.AnnualRenewal === value) {
    return 'annually';
  }
  if (ContractRenewalCycle.MonthlyRenewal === value) {
    return 'monthly';
  }
}

export const ContractCard = ({
  data,
  organizationName,
  organizationId,
}: ContractCardProps) => {
  const queryKey = useGetContractsQuery.getKey({ id: organizationId });
  const queryClient = useQueryClient();

  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isExpanded, setIsExpanded] = useState(!data?.signedAt);
  const formId = 'contractForm';
  const client = getGraphQLClient();

  const updateContract = useUpdateContractMutation(client, {
    // todo fix https://linear.app/customer-os/issue/COS-985/fix-optimitsic-update-for-update-contract-mutation
    // onMutate: ({ input }) => {
    //   queryClient.cancelQueries({ queryKey });
    //
    //   queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
    //     return produce(currentCache, (draft) => {
    //       const previousContracts = draft?.['organization']?.['contracts'];
    //       const updatedContractIndex = previousContracts?.findIndex(
    //         (contract) => contract.id === data?.id,
    //       );
    //       if (draft?.['organization']?.['contracts']) {
    //         draft['organization']['contracts']?.map((contractData, index) => {
    //           if (index !== updatedContractIndex) {
    //             return contractData;
    //           }
    //
    //           return { ...contractData, ...input };
    //         });
    //       }
    //     });
    //   });
    //   const previousEntries =
    //     queryClient.getQueryData<GetContractsQuery>(queryKey);
    //
    //   return { previousEntries };
    // },
    // onError: (_, __, context) => {
    //   queryClient.setQueryData<GetContractsQuery>(
    //     queryKey,
    //     context?.previousEntries,
    //   );
    //   toastError('Failed to update contract', 'update-contract-error');
    // },
    onSettled: () => {
      queryClient.invalidateQueries(queryKey);
    },
  });

  const updateContractDebounced = debounce(
    (variables: { input: ContractUpdateInput }) => {
      updateContract.mutate({
        ...variables,
      });
    },
    300,
  );

  const defaultValues = ContractDTO.toForm(data);
  const { setDefaultValues } = useForm<TimeToRenewalForm>({
    formId,
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'name') {
          updateContractDebounced({
            input: {
              contractId: data.id,
              ...ContractDTO.toPayload({
                ...state.values,
                [action.payload.name]: action.payload.value,
              }),
            },
          });

          return next;
        }
        if (action.payload.name === 'contractUrl') {
          return next;
        }
        updateContract.mutate({
          input: {
            contractId: data.id,
            ...ContractDTO.toPayload({
              ...state.values,
              [action.payload.name]: action.payload.value,
            }),
          },
        });
      }

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
        pb={isExpanded ? 2 : 0}
        w='full'
        flexDir='column'
      >
        <Flex justifyContent='space-between' w='full' flex={1}>
          <Heading
            size='sm'
            color='gray.700'
            noOfLines={1}
            alignItems='baseline'
            as={Flex}
          >
            {!isExpanded && (data?.name || `${organizationName} contract`)}

            {isExpanded && (
              <FormInput
                fontWeight='semibold'
                fontSize='inherit'
                height='fit-content'
                name='name'
                formId={formId}
                placeholder='Contract name'
                borderBottom='none'
                _hover={{
                  borderBottom: 'none',
                }}
              />
            )}
          </Heading>
          <Flex alignItems='center' gap={2} ml={4}>
            <UrlInput
              formId={formId}
              url={data?.contractUrl}
              contractId={data?.id}
              onSubmit={updateContract.mutate}
            />

            <ContractStatusSelect status={data.status} />

            {isExpanded && (
              <IconButton
                size='xs'
                variant='ghost'
                aria-label='Collapse'
                icon={<ChevronUp />}
                onClick={() => setIsExpanded(false)}
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
            height='fit-content'
            alignItems='flex-start'
            justifyContent='flex-start'
            onClick={() => setIsExpanded(true)}
          >
            <Flex
              flexDir='column'
              alignItems='flex-start'
              justifyContent='center'
            >
              {!data?.signedAt && <Text>No start date or services yet</Text>}

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
                      data.status !== ContractStatus.Ended &&
                      `Renews ${getLabelFromValue(
                        data.renewalCycle,
                      )} on ${DateTimeUtils.format(
                        calculateNextRenewalDate(
                          data.serviceStartedAt,
                          data.renewalCycle,
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

            <Edit03 ml={1} color='gray.400' boxSize='3' mt='3px' />
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
              placeholder='Renewal cycle'
              isLabelVisible
              name='renewalCycle'
              formId={formId}
              options={billingFrequencyOptions}
            />
          </Flex>
        </CardBody>
      )}
      <CardFooter p='0' mt={1} w='full' flexDir='column'>
        {data?.serviceStartedAt && data?.renewalCycle && (
          <RenewalARRCard
            hasEnded={data.status === ContractStatus.Ended}
            startedAt={data.serviceStartedAt}
            renewCycle={data.renewalCycle}
          />
        )}
        <Services contractId={data.id} data={data?.serviceLineItems} />
      </CardFooter>
    </Card>
  );
};

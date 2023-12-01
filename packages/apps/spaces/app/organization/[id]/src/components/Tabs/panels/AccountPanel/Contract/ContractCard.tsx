import { useForm } from 'react-inverted-form';
import React, { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
import { debounce } from 'lodash';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { Check } from '@ui/media/icons/Check';
import { Edit03 } from '@ui/media/icons/Edit03';
import { FormSelect } from '@ui/form/SyncSelect';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { toastError } from '@ui/presentation/Toast';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardBody, CardFooter, CardHeader } from '@ui/presentation/Card';
import { Contract, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { useUpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { ContractSubtitle } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractSubtitle';

import { UrlInput } from './UrlInput';
import { Services } from './Services/Services';
import { billingFrequencyOptions } from '../utils';
import { RenewalARRCard } from './RenewalARR/RenewalARRCard';
import { ContractDTO, TimeToRenewalForm } from './Contract.dto';
import { ContractStatusSelect } from './contractStatuses/ContractStatusSelect';

interface ContractCardProps {
  data: Contract;
  organizationId: string;
  organizationName: string;
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
  const formId = `contract-form-${data.id}`;

  const client = getGraphQLClient();
  const updateContract = useUpdateContractMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          const updatedContractIndex = previousContracts?.findIndex(
            (contract) => contract.id === data?.id,
          );
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              if (index !== updatedContractIndex) {
                return contractData;
              }

              return { ...input };
            });
          }
        });
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },
    onError: (error, { input }, context) => {
      queryClient.setQueryData<GetContractsQuery>(
        queryKey,
        context?.previousEntries,
      );

      const invalidDate =
        DateTimeUtils.isBefore(input.endedAt, input.serviceStartedAt) ||
        DateTimeUtils.isBefore(input.endedAt, input.signedAt);

      toastError(
        `${
          invalidDate
            ? 'The contract end date needs to be after the service start date'
            : 'Failed to update contract'
        }`,
        `update-contract-error-${error}`,
      );
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
      }, 1000);
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

  const defaultValues = ContractDTO.toForm({
    organizationName,
    ...(data ?? {}),
  });
  const { setDefaultValues, state } = useForm<TimeToRenewalForm>({
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
      updateContractDebounced.flush();
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Card
      as='section'
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
        role='button'
        pb={isExpanded ? 2 : 0}
        w='full'
        flexDir='column'
        _hover={
          !isExpanded
            ? {
                '#edit-contract-icon': {
                  opacity: 1,
                  transition: 'opacity 0.2s linear',
                },
              }
            : {}
        }
        sx={
          !isExpanded
            ? {
                '#edit-contract-icon': {
                  opacity: 0,
                  transition: 'opacity 0.2s linear',
                },
              }
            : {}
        }
        onClick={() => (!isExpanded ? setIsExpanded(true) : null)}
      >
        <Flex justifyContent='space-between' w='full' flex={1}>
          <Heading
            size='sm'
            color='gray.700'
            noOfLines={1}
            lineHeight={1.4}
            display='inline'
            w={isExpanded ? '235px' : '250px'}
            whiteSpace='nowrap'
          >
            {!isExpanded && state.values.name}

            {isExpanded && (
              <FormInput
                fontWeight='semibold'
                fontSize='inherit'
                height='fit-content'
                name='name'
                formId={formId}
                borderBottom='none'
                _hover={{
                  borderBottom: 'none',
                }}
              />
            )}
          </Heading>

          <Flex alignItems='center' gap={2} ml={2}>
            {!isExpanded && (
              <Edit03
                mr={1}
                color='gray.400'
                boxSize='4'
                id='edit-contract-icon'
              />
            )}
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
                icon={<Check color='gray.400' />}
                onClick={() => setIsExpanded(false)}
              />
            )}
          </Flex>
        </Flex>

        {!isExpanded && (
          <Flex
            bg='transparent'
            _hover={{
              bg: 'transparent',
              svg: { opacity: 1, transition: 'opacity 0.2s linear' },
            }}
            sx={{ svg: { opacity: 0, transition: 'opacity 0.2s linear' } }}
            fontSize='sm'
            fontWeight='normal'
            color='gray.500'
            p={0}
            height='fit-content'
            alignItems='flex-start'
            justifyContent='flex-start'
          >
            <ContractSubtitle data={data} />
          </Flex>
        )}
      </CardHeader>
      {isExpanded && (
        <CardBody as={Flex} p='0' flexDir='column' w='full'>
          <Flex gap='4' mb={2} flexGrow={0}>
            <DatePicker
              label='Contract signed'
              placeholder='Signed date'
              formId={formId}
              name='signedAt'
              inset='120% auto auto 0px'
              calendarIconHidden
            />
            <DatePicker
              label='Contract ends'
              placeholder='End date'
              minDate={state.values.serviceStartedAt}
              formId={formId}
              name='endedAt'
              calendarIconHidden
            />
          </Flex>
          <Flex gap='4' flexGrow={0}>
            <DatePicker
              label='Service starts'
              placeholder='Start date'
              formId={formId}
              name='serviceStartedAt'
              inset='120% auto auto 0px'
              maxDate={state.values.endedAt}
              calendarIconHidden
            />
            <FormSelect
              label='Contract renews'
              placeholder='Renewal cycle'
              isLabelVisible
              name='renewalCycle'
              formId={formId}
              options={billingFrequencyOptions}
              // isClearable
            />
          </Flex>
        </CardBody>
      )}
      <CardFooter p='0' mt={1} w='full' flexDir='column'>
        {data?.opportunities && data.renewalCycle && (
          <RenewalARRCard
            hasEnded={data.status === ContractStatus.Ended}
            startedAt={data.serviceStartedAt}
            renewCycle={data.renewalCycle}
            opportunity={data.opportunities?.[0]}
          />
        )}
        <Services contractId={data.id} data={data?.serviceLineItems} />
      </CardFooter>
    </Card>
  );
};

import { useRef, useEffect, useState } from 'react';
import { useForm } from 'react-inverted-form';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { BillingDetails } from '@graphql/types';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { Divider } from '@ui/presentation/Divider';
import { FormSelect } from '@ui/form/SyncSelect';
import { DateTimeUtils } from '@spaces/utils/date';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { invalidateAccountDetailsQuery } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { useQueryClient } from '@tanstack/react-query';
import { useUpdateBillingDetailsMutation } from '@organization/src/graphql/updateBillingDetails.generated';

import { getTimeToRenewal } from '../../../shared/util';
import { frequencyOptions } from '../utils';
import { TimeToRenewalDTO } from './TimeToRenewal.dto';

interface TimeToRenewalsCardProps {
  id: string;
  data?: BillingDetails | null;
}
export const TimeToRenewal = ({ id, data }: TimeToRenewalsCardProps) => {
  const [isFocused, setIsFocused] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryClient = useQueryClient();
  const client = getGraphQLClient();
  const updateBillingDetails = useUpdateBillingDetailsMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
    },
  });

  const defaultValues = TimeToRenewalDTO.toForm(data);

  useForm({
    formId: 'time-to-renewal',
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        setIsFocused(false);
      }
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'renewalCycle': {
            const renewalCycle = action.payload?.value;
            const renewalCycleStart = state.values.renewalCycleStart;

            if (!renewalCycle && renewalCycleStart !== null) {
              updateBillingDetails.mutate({
                input: {
                  ...TimeToRenewalDTO.toPayload({
                    id,
                    ...data,
                    renewalCycle: null,
                    renewalCycleStart: null,
                  }),
                },
              });
              return {
                ...next,
                values: {
                  ...next.values,
                  renewalCycleStart: null,
                },
              };
            }
            updateBillingDetails.mutate({
              input: {
                ...TimeToRenewalDTO.toPayload({
                  id,
                  ...data,
                  renewalCycle,
                  renewalCycleStart,
                }),
              },
            });
            break;
          }
          case 'renewalCycleStart': {
            updateBillingDetails.mutate({
              input: {
                ...TimeToRenewalDTO.toPayload({
                  id,
                  ...data,
                  ...state.values,
                  renewalCycleStart: action.payload?.value || null,
                }),
              },
            });
            break;
          }
        }
      }

      return next;
    },
  });

  const [value, label] = getTimeToRenewal(data?.renewalCycleNext);

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Card
      p='4'
      w='full'
      size='lg'
      variant='outline'
      cursor='default'
      boxShadow={isFocused ? 'md' : 'xs'}
      _hover={{
        boxShadow: 'md',
      }}
      transition='all 0.2s ease-out'
    >
      <CardBody as={Flex} p='0' justify='space-between' align='center' w='full'>
        <FeaturedIcon size='md' minW='10'>
          <Icons.ClockFastForward />
        </FeaturedIcon>
        <Flex ml='5' align='center' justify='space-between' w='full'>
          <Flex flexDir='column'>
            <Heading size='sm' color='gray.700'>
              Time to renewal
            </Heading>
            <Text fontSize='xs' color='gray.500'>
              {data?.renewalCycleNext
                ? `Renews on ${DateTimeUtils.format(
                    data?.renewalCycleNext,
                    DateTimeUtils.dateWithFullMonth,
                  )}`
                : 'Add a renewal cycle and start date'}
            </Text>
          </Flex>

          <Flex direction='column' alignItems='flex-end' justifyItems='center'>
            <Text
              fontSize='2xl'
              fontWeight='bold'
              lineHeight='1'
              color={!data?.renewalCycleNext ? 'gray.400' : 'gray.700'}
            >
              {data?.renewalCycleNext ? value : 'Unknown'}
            </Text>
            {data?.renewalCycleNext && <Text color='gray.500'>{label}</Text>}
          </Flex>
        </Flex>
      </CardBody>

      <CardFooter as={Flex} p='0' w='full' flexDir='column'>
        <Divider my='4' />
        <Flex gap='4'>
          <FormSelect
            isClearable
            label='Renewal cycle'
            isLabelVisible
            name='renewalCycle'
            placeholder='Monthly'
            options={frequencyOptions}
            formId='time-to-renewal'
            onFocus={() => setIsFocused(true)}
            leftElement={<Icons.ClockFastForward mr='3' color='gray.500' />}
          />
          <DatePicker
            label='Renewal cycle start'
            formId='time-to-renewal'
            name='renewalCycleStart'
            onFocus={() => setIsFocused(true)}
          />
        </Flex>
      </CardFooter>
    </Card>
  );
};

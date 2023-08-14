import React from 'react';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';
import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Divider } from '@ui/presentation/Divider';
import { VStack } from '@ui/layout/Stack';
import { Heading } from '@ui/typography/Heading';
import BillingDetails from '@spaces/atoms/icons/BillingDetails';
import CurrencyDollar from '@spaces/atoms/icons/CurrencyDollar';
import { FormSelect } from '@ui/form/SyncSelect';
import CoinsSwap from '@spaces/atoms/icons/CoinsSwap';
import { frequencyOptions } from './utils';
import ClockCheck from '@spaces/atoms/icons/ClockCheck';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { FormCurrencyInput } from '@ui/form/CurrencyInput/FormCurrencyInput';
import { useForm } from 'react-inverted-form';
import {
  OrganizationAccountBillingDetails,
  OrganizationAccountBillingDetailsForm,
} from './OrganziationAccountBillingDetails.dto';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { invalidateAccountDetailsQuery } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { useQueryClient } from '@tanstack/react-query';
import { useUpdateBillingDetailsMutation } from '@organization/graphql/updateBillingDetails.generated';
import { BillingDetails as BillingDetailsT } from '@graphql/types';

export type BillingDetailsType =  BillingDetailsT & { amount?: string | null }
interface BillingDetailsCardBProps {
  billingDetailsData: BillingDetailsType;
  id: string;
}
export const BillingDetailsCard: React.FC<BillingDetailsCardBProps> = ({
  billingDetailsData,
  id,
}) => {
  const queryClient = useQueryClient();
  const defaultValues: OrganizationAccountBillingDetailsForm =
    new OrganizationAccountBillingDetails(billingDetailsData);
  const client = getGraphQLClient();
  const updateBillingDetails = useUpdateBillingDetailsMutation(client, {
    onSuccess: () => invalidateAccountDetailsQuery(queryClient, id),
  });
  const handleUpdateBillingDetails = (
    variables: Partial<OrganizationAccountBillingDetailsForm>,
  ) => {
    const inputData = OrganizationAccountBillingDetails.toPayload({
      ...state.values,
      ...variables,
    });
    updateBillingDetails.mutate({
      input: { id, ...inputData },
    });
  };

  const formId = 'organization-account-billing-details-form';
  const { state } = useForm<OrganizationAccountBillingDetailsForm>({
    formId,
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        const shouldPreventFrequencyOptionSave =
          action.payload?.value?.value ===
          //@ts-expect-error fixme
          defaultValues?.[action.payload.name]?.value;

        switch (action.payload.name) {
          case 'frequency': {
            if (shouldPreventFrequencyOptionSave) {
              return next;
            }
            handleUpdateBillingDetails({
              frequency: action.payload?.value?.value || null,
            });

            return next;
          }
          case 'renewalCycle': {
            if (shouldPreventFrequencyOptionSave) {
              return next;
            }
            const renewalCycle = action.payload?.value?.value;
            const renewalCycleStart = state.values.renewalCycleStart;

            if (!renewalCycle && renewalCycleStart !== null) {
              handleUpdateBillingDetails({
                renewalCycle: null,
                renewalCycleStart: null,
              });
              return {
                ...next,
                values: {
                  ...next.values,
                  renewalCycleStart: null,
                },
              };
            }
            handleUpdateBillingDetails({
              renewalCycle,
            });
            return next;
          }
          case 'renewalCycleStart': {
            const shouldPreventSave =
              //@ts-expect-error fixme
              action.payload?.value === defaultValues?.[action.payload.name];
            if (shouldPreventSave) return next;
            handleUpdateBillingDetails({
              renewalCycleStart: action.payload?.value || null,
            });
            return next;
          }
          default:
            return next;
        }
      }

      if (action.type === 'FIELD_BLUR' && action.payload.name === 'amount') {
        if (defaultValues.amount === action.payload.value) return next;
        handleUpdateBillingDetails({
          amount: action.payload.value,
        });
      }

      return next;
    },
  });

  return (
    <Card
      size='sm'
      width='full'
      borderRadius='xl'
      border='1px solid'
      borderColor='gray.200'
      boxShadow='xs'
    >
      <CardHeader display='flex' alignItems='center'>
        <BillingDetails />
        <Heading ml={5} size='sm' color='gray.700'>
          Billing details
        </Heading>
      </CardHeader>
      <Box px={4}>
        <Divider color='gray.200' />
      </Box>

      <CardBody padding={4}>
        <VStack spacing='4' w='full'>
          <Flex justifyItems='space-between' w='full'>
            <FormCurrencyInput
              label='Billing amounts'
              color='gray.700'
              isLabelVisible
              formId={formId}
              name='amount'
              min={0}
              placeholder='Amount'
              leftElement={
                <Box color='gray.500'>
                  <CurrencyDollar height='16px' />
                </Box>
              }
            />

            <FormSelect
              isClearable
              label='Billing frequency'
              isLabelVisible
              name='frequency'
              placeholder='Monthly'
              options={frequencyOptions}
              formId={formId}
              leftElement={
                <Box mr={3} color='gray.500'>
                  <CoinsSwap height={16} />
                </Box>
              }
            />
          </Flex>
          <Flex justifyItems='space-between' w='full'>
            <FormSelect
              isClearable
              label='Renewal cycle'
              isLabelVisible
              name='renewalCycle'
              placeholder='Monthly'
              options={frequencyOptions}
              formId={formId}
              leftElement={
                <Box mr={3} color='gray.500'>
                  <ClockCheck height={16} />
                </Box>
              }
            />
            <DatePicker
              label='Renewal cycle start'
              formId={formId}
              name='renewalCycleStart'
            />
          </Flex>
        </VStack>
      </CardBody>
    </Card>
  );
};

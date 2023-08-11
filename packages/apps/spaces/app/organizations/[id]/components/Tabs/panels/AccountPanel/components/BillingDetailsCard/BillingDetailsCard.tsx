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
import { frequencyOptions } from '@organization/components/Tabs/panels/AccountPanel/components/BillingDetailsCard/utils';
import ClockCheck from '@spaces/atoms/icons/ClockCheck';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { useForm } from 'react-inverted-form';

import {
  OrganizationAccountBillingDetails,
  OrganizationAccountBillingDetailsForm,
} from '@organization/components/Tabs/panels/AccountPanel/components/BillingDetailsCard/OrganziationAccountBillingDetails.dto';
import { FormCurrencyInput } from '@ui/form/CurrencyInput/FormCurrencyInput';

export const BillingDetailsCard: React.FC = () => {
  const defaultValues: OrganizationAccountBillingDetailsForm =
    new OrganizationAccountBillingDetails();
  const formId = 'organization-account-form';
  const { state } = useForm<OrganizationAccountBillingDetailsForm>({
    formId,
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        // TODO uncomment when BE is connected
        // const shouldPreventSave =
        //   action.payload?.value?.value ===
        //   //@ts-expect-error fixme
        //   defaultValues?.[action.payload.name]?.value;
        // if (shouldPreventSave) {
        //   return next;
        // }
        switch (action.payload.name) {
          case 'billingDetailsRenewalCycle': {
            const renewalCycle = action.payload?.value?.value;
            const renewalCycleStart =
              state.values.billingDetailsRenewalCycleStart;

            if (!renewalCycle && renewalCycleStart !== null) {
              return {
                ...next,
                values: {
                  ...next.values,
                  billingDetailsRenewalCycleStart: null,
                },
              };
            }

            return {
              ...next,
              values: {
                ...next.values,
                stage: null,
              },
            };
          }
          default:
            return next;
        }
      }

      if (action.type === 'FIELD_BLUR') {
        switch (action.payload.name) {
          case 'billingDetailsAmount': {
            const trimmedValue = (action.payload?.value || '')?.trim();
            if (
              //@ts-expect-error fixme
              state.fields?.[action.payload.name].meta.pristine ||
              //@ts-expect-error fixme
              trimmedValue === defaultValues?.[action.payload.name]
            ) {
              return next;
            }
            break;
          }
          default:
            return next;
        }
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
              name='billingDetailsAmount'
              min={0}
              placeholder='$1700'
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
              name='billingDetailsFrequency'
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
              name='billingDetailsRenewalCycle'
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
              name='billingDetailsRenewalCycleStart'
            />
          </Flex>
        </VStack>
      </CardBody>
    </Card>
  );
};

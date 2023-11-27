'use client';
import React, { useRef } from 'react';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { FormNumberInput } from '@ui/form/NumberInput';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { FormCurrencyInput } from '@ui/form/CurrencyInput';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

interface RecurringServiceFormProps {
  formId: string;
}

const [_, _1, ...subscriptionOptions] = billedTypeOptions;
export const RecurringServiceFrom = ({ formId }: RecurringServiceFormProps) => {
  const initialRef = useRef(null);

  return (
    <>
      <FormInput
        name='name'
        formId={formId}
        label='Service name'
        placeholder="What's this service's name?"
        isLabelVisible
        labelProps={{
          fontSize: 'sm',
          fontWeight: 'semibold',
          mb: 0,
          lineHeight: 1,
        }}
      />
      <Flex gap={4} mt={2} justifyContent='space-between'>
        <FormNumberInput
          w='full'
          height='auto'
          placeholder='Quantity'
          isLabelVisible
          label='Licences'
          min={1}
          ref={initialRef}
          leftElement={<Certificate02 boxSize={4} color='gray.500' />}
          formId={formId}
          name='quantity'
        />

        <FormCurrencyInput
          formId={formId}
          name='price'
          w='full'
          placeholder='Per license'
          isLabelVisible
          label='Price/qty'
          min={1}
          leftElement={<CurrencyDollar boxSize={4} color='gray.500' />}
        />

        <FormSelect
          label='Billed'
          placeholder='Frequency'
          isLabelVisible
          name='billed'
          formId={formId}
          options={subscriptionOptions}
          leftElement={<ClockCheck mr='3' color='gray.500' boxSize={4} />}
        />
      </Flex>
    </>
  );
};

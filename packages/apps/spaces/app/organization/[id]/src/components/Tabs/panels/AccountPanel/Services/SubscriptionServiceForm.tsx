'use client';
import React, { useRef, useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { FormNumberInput } from '@ui/form/NumberInput';
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

interface SubscriptionServiceFromProps {
  data: {
    name?: string | null;
    licenses?: string | null;
    licensePrice?: string | null;
  };
}

export const SubscriptionServiceFrom = ({
  data,
}: SubscriptionServiceFromProps) => {
  const initialRef = useRef(null);
  const formId = 'TODO';

  const [licensePrice, setLicensePrice] = useState<string>(
    data?.licensePrice || '',
  );
  const [licenses, setLicenses] = useState<string>(data?.licenses || '');
  const [name, setName] = useState<string>(data?.name || '');

  return (
    <>
      <FormInput
        name='name'
        formId='todo'
        value={name}
        onChange={(e) => setName(e.target.value)}
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
          onChange={setLicenses}
          value={`${licenses}`}
          w='full'
          height='auto'
          placeholder='Quantity'
          isLabelVisible
          label='Licences'
          min={0}
          ref={initialRef}
          leftElement={<Certificate02 boxSize={4} color='gray.500' />}
          formId={formId}
          name='licences'
        />

        <CurrencyInput
          onChange={setLicensePrice}
          value={`${licensePrice}`}
          w='full'
          placeholder='Per license'
          isLabelVisible
          label='Price/license'
          min={0}
          ref={initialRef}
          leftElement={<CurrencyDollar boxSize={4} color='gray.500' />}
        />

        <FormSelect
          label='Billed'
          placeholder='Frequency'
          isLabelVisible
          name='billingFrequency'
          formId={formId}
          options={billingFrequencyOptions}
          leftElement={<ClockCheck mr='3' color='gray.500' boxSize={4} />}
        />
      </Flex>
    </>
  );
};

'use client';
import { useRef, useState } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { RenewalCycle } from '@graphql/types';
import { FormSelect } from '@ui/form/SyncSelect';
import { FormNumberInput } from '@ui/form/NumberInput';
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { SelectOption } from '@shared/types/SelectOptions';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { frequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

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
  const [billingFrequency, setBillingFrequency] = useState<
    SelectOption<RenewalCycle>
  >(frequencyOptions[2]);
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
        label='Service Name'
        placeholder='What’s this service about?'
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
          placeholder='Quantity'
          isLabelVisible
          label='Licences'
          min={0}
          ref={initialRef}
          leftElement={
            <Box color='gray.500'>
              <Certificate02 height='16px' />
            </Box>
          }
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
          leftElement={
            <Box color='gray.500'>
              <CurrencyDollar height='16px' />
            </Box>
          }
        />
        <FormSelect
          label='Billed'
          isLabelVisible
          name='billingFrequency'
          formId='tbd'
          value={billingFrequency}
          onChange={(d) => setBillingFrequency(d)}
          options={frequencyOptions}
          leftElement={<ClockCheck mr='3' color='gray.500' />}
        />
      </Flex>
    </>
  );
};

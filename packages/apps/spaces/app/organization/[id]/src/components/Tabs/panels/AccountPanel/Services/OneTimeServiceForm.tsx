'use client';
import { useRef, useState } from 'react';

import { Box } from '@ui/layout/Box';
import { FormInput } from '@ui/form/Input';
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';

export type OneTimeServiceValue = {
  name?: string | null;
  price?: string | null;
};

interface OneTimeServiceModalProps {
  data: OneTimeServiceValue;
}

export const OneTimeServiceForm = ({ data }: OneTimeServiceModalProps) => {
  const initialRef = useRef(null);

  const [price, setPrice] = useState<string>(data?.price || '');
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
        mb={2}
        labelProps={{
          fontSize: 'sm',
          fontWeight: 'semibold',
          mb: 0,
          lineHeight: 1,
        }}
      />
      <CurrencyInput
        onChange={setPrice}
        value={`${price}`}
        w='full'
        placeholder='Price'
        isLabelVisible
        label='Price'
        min={0}
        ref={initialRef}
        leftElement={
          <Box color='gray.500'>
            <CurrencyDollar height='16px' />
          </Box>
        }
      />
    </>
  );
};

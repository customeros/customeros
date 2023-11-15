'use client';
import { useRef, useState } from 'react';

import { CurrencyInput } from '@ui/form/CurrencyInput';
import { Input, FormLabel, FormControl } from '@ui/form/Input';
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
      <FormControl>
        <FormLabel
          fontWeight='semibold'
          color='gray.700'
          fontSize='sm'
          mb={0}
          lineHeight={1}
        >
          Service name
        </FormLabel>

        <Input
          name='name'
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="What's this service's name?"
          mb={2}
          autoComplete='off'
        />
      </FormControl>

      <CurrencyInput
        onChange={setPrice}
        value={`${price}`}
        w='full'
        height='auto'
        placeholder='Price'
        isLabelVisible
        label='Price'
        min={0}
        ref={initialRef}
        leftElement={<CurrencyDollar boxSize={4} color='gray.500' />}
      />
    </>
  );
};

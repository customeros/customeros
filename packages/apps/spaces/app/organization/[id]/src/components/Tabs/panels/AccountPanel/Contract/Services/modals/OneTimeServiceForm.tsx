'use client';

import { BilledType } from '@graphql/types';
import { FormInput, FormControl } from '@ui/form/Input';
import { FormCurrencyInput } from '@ui/form/CurrencyInput';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';

interface OneTimeServiceModalProps {
  formId: string;
  billedType: BilledType;
}

export const OneTimeServiceForm = ({
  formId,
  billedType,
}: OneTimeServiceModalProps) => {
  return (
    <>
      <FormControl>
        <FormInput
          formId={formId}
          name='name'
          label='Service name'
          placeholder="What's this service's name?"
          mb={2}
          autoComplete='off'
          isLabelVisible
          labelProps={{
            fontSize: 'sm',
            fontWeight: 'semibold',
            mb: 0,
            lineHeight: 1,
          }}
        />
      </FormControl>

      <FormCurrencyInput
        name='price'
        formId={formId}
        w='full'
        height='auto'
        placeholder={
          billedType === BilledType.Once ? 'One-time price' : 'Price'
        }
        isLabelVisible
        label={billedType === BilledType.Once ? 'Price' : 'Price/usage'}
        min={0}
        leftElement={<CurrencyDollar boxSize={4} color='gray.500' />}
      />
    </>
  );
};

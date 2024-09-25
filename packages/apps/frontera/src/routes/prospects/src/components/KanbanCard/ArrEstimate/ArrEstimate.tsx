import { useState, useEffect } from 'react';

import { P, match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Currency } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { MaskedInput } from '@ui/form/Input/MaskedInput';
import { currencySymbol } from '@shared/util/currencyOptions';

interface ArrEstimateProps {
  opportunityId: string;
}

export const ArrEstimate = observer(({ opportunityId }: ArrEstimateProps) => {
  const store = useStore();
  const opportunity = store.opportunities.value.get(opportunityId);
  const defaultValue = opportunity?.value?.maxAmount ?? 0;
  const [value, setValue] = useState(defaultValue.toString());

  const handleAccept = (unmaskedValue: string) => {
    opportunity?.update(
      (value) => {
        value.maxAmount = unmaskedValue ? parseFloat(unmaskedValue) : 0;

        return value;
      },
      { mutate: false },
    );
  };

  const handleBlur = () => {
    opportunity?.saveProperty('maxAmount');
  };

  const defaultCurrency = match(store.settings.tenant.value?.baseCurrency)
    .with(P.nullish, () => Currency.Usd)
    .with(P.string, (str) => (str.length === 3 ? str : Currency.Usd))
    .otherwise((tenantCurrency) => tenantCurrency);

  const symbol = match(opportunity?.value?.currency)
    .with(P.nullish, () => currencySymbol[defaultCurrency])
    .otherwise(
      (currency) => currencySymbol[currency] ?? currencySymbol[defaultCurrency],
    );

  useEffect(() => {
    setValue(defaultValue.toString());
  }, [defaultValue]);

  return (
    <MaskedInput
      size='xs'
      value={value}
      maxLength={12}
      variant='unstyled'
      onBlur={handleBlur}
      mask={`${symbol}num`}
      className='max-w-[100px]'
      placeholder='ARR estimate'
      defaultValue={defaultValue.toString()}
      onClick={(e) => (e.target as HTMLInputElement).select()}
      onAccept={(v, instance) => {
        setValue(v);
        handleAccept(instance._unmaskedValue);
      }}
      blocks={{
        num: {
          mask: Number,
          scale: 0,
          lazy: false,
          min: 0,
          placeholderChar: '#',
          thousandsSeparator: ',',
          normalizeZeros: true,
          padFractionalZeros: true,
          autofix: true,
        },
      }}
    />
  );
});

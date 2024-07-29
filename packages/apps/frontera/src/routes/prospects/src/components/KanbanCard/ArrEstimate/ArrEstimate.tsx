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

  const symbol =
    currencySymbol[store.settings.tenant.value?.baseCurrency ?? Currency.Usd];

  return (
    <MaskedInput
      size='xs'
      variant='unstyled'
      onBlur={handleBlur}
      mask={`${symbol}num`}
      placeholder='ARR estimate'
      defaultValue={defaultValue.toString()}
      onAccept={(_, instance) => handleAccept(instance._unmaskedValue)}
      blocks={{
        num: {
          mask: Number,
          scale: 2,
          lazy: false,
          placeholderChar: '#',
          thousandsSeparator: ',',
          normalizeZeros: true,
          padFractionalZeros: true,
          radix: '.',
          autofix: true,
        },
      }}
    />
  );
});

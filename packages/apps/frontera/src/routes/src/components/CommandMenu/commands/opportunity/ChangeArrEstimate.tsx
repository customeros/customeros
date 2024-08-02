import { useState } from 'react';

import { P, match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Currency } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { MaskedInput } from '@ui/form/Input/MaskedInput';
import { currencySymbol } from '@shared/util/currencyOptions';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeArrEstimate = observer(() => {
  const store = useStore();
  const [value, setValue] = useState('');
  const [unmaskedValue, setUnmaskedValue] = useState('');
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(
    (context.ids as string[])?.[0],
  );

  const label = match(context.entity)
    .with('Opportunity', () => `Opportunity - ${opportunity?.value?.name}`)
    .otherwise(() => 'Change ARR estimate');

  const defaultCurrency = match(store.settings.tenant.value?.baseCurrency)
    .with(P.nullish, () => Currency.Usd)
    .with(P.string, (str) => (str.length === 3 ? str : Currency.Usd))
    .otherwise((tenantCurrency) => tenantCurrency);

  const symbol = match(opportunity?.value?.currency)
    .with(P.nullish, () => currencySymbol[defaultCurrency])
    .otherwise(
      (currency) => currencySymbol[currency] ?? currencySymbol[defaultCurrency],
    );

  const handleEnterKey = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      opportunity?.update((value) => {
        value.maxAmount = parseFloat(unmaskedValue);

        return value;
      });
      store.ui.commandMenu.setType('OpportunityCommands');
      store.ui.commandMenu.setOpen(false);
    }
  };

  return (
    <Command shouldFilter={false} onKeyDown={handleEnterKey}>
      <CommandInput
        asChild
        label={label}
        placeholder='Type a command or search'
      >
        <MaskedInput
          size='xs'
          value={value}
          variant='unstyled'
          mask={`${symbol}num`}
          placeholder='Change ARR estimate'
          onAccept={(val, maskRef) => {
            setValue(val);
            setUnmaskedValue(maskRef?.unmaskedValue);
          }}
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
      </CommandInput>

      <Command.List className='p-0'></Command.List>
    </Command>
  );
});

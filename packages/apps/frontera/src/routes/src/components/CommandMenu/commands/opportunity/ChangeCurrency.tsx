import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Currency } from '@graphql/types';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { CurrencyEuro } from '@ui/media/icons/CurrencyEuro';
import { CurrencyPound } from '@ui/media/icons/CurrencyPound';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeCurrency = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(
    (context.ids as string[])?.[0],
  );

  const label = match(context.entity)
    .with('Opportunity', () => `Opportunity - ${opportunity?.value?.name}`)
    .otherwise(() => undefined);

  const handleSelect = (currency: Currency) => {
    opportunity?.update((value) => {
      Object.assign(value, { currency });

      return value;
    });

    store.ui.commandMenu.setOpen(false);
  };

  return (
    <Command label='Change currency'>
      <CommandInput label={label} placeholder='Change ARR currency...' />

      <Command.List>
        <CommandItem
          leftAccessory={<CurrencyDollar />}
          onSelect={() => handleSelect(Currency.Usd)}
          rightAccessory={
            opportunity?.value?.currency === Currency.Usd && <Check />
          }
        >
          USD
        </CommandItem>
        <CommandItem
          leftAccessory={<CurrencyEuro />}
          onSelect={() => handleSelect(Currency.Eur)}
          rightAccessory={
            opportunity?.value?.currency === Currency.Eur && <Check />
          }
        >
          EUR
        </CommandItem>
        <CommandItem
          leftAccessory={<CurrencyPound />}
          onSelect={() => handleSelect(Currency.Gbp)}
          rightAccessory={
            opportunity?.value?.currency === Currency.Gbp && <Check />
          }
        >
          GBP
        </CommandItem>
      </Command.List>
    </Command>
  );
});

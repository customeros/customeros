import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { CommandMenuType } from '@store/UI/CommandMenu.store';

import { Currency } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { CurrencyEuro } from '@ui/media/icons/CurrencyEuro';
import { CurrencyPound } from '@ui/media/icons/CurrencyPound';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeCurrency = observer(() => {
  const store = useStore();
  const location = useLocation();
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(context.id as string);

  const label = match(context.entity)
    .with('Opportunity', () => `Opportunity - ${opportunity?.value?.name}`)
    .otherwise(() => undefined);

  const hubType = match(location.pathname)
    .with('/prospects', () => 'OpportunityHub')
    .otherwise(() => 'GlobalHub') as CommandMenuType;

  const handleSelect = (currency: Currency) => {
    opportunity?.update((value) => {
      Object.assign(value, { currency });

      return value;
    });

    store.ui.commandMenu.setOpen(false);
  };

  useEffect(() => {
    () => {
      return match(context.entity)
        .with('Opportunity', () => {
          store.ui.commandMenu.setType(hubType);
        })
        .otherwise(() => {});
    };
  }, []);

  return (
    <Command label='Change currency'>
      <CommandInput label={label} placeholder='Change ARR currency...' />

      <Command.List>
        <CommandItem
          leftAccessory={<CurrencyDollar />}
          onSelect={() => handleSelect(Currency.Usd)}
        >
          USD
        </CommandItem>
        <CommandItem
          leftAccessory={<CurrencyEuro />}
          onSelect={() => handleSelect(Currency.Eur)}
        >
          EUR
        </CommandItem>
        <CommandItem
          leftAccessory={<CurrencyPound />}
          onSelect={() => handleSelect(Currency.Gbp)}
        >
          GBP
        </CommandItem>
      </Command.List>
    </Command>
  );
});
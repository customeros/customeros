import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Currency } from '@graphql/types';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { CurrencyEuro } from '@ui/media/icons/CurrencyEuro';
import { CurrencyPound } from '@ui/media/icons/CurrencyPound';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';

export const AddCurrencySubItemGroup = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const opportunity = store.opportunities.value.get(
    (context.ids as string[])?.[0],
  );

  const handleSelect = (currency: Currency) => {
    match(context.entity)
      .with('Opportunity', () => {
        opportunity?.update((value) => {
          Object.assign(value, { currency });

          return value;
        });
      })
      .with('Opportunities', () => {
        context.ids?.forEach((id) => {
          const opportunity = store.opportunities.value.get(id);

          opportunity?.update((value) => {
            Object.assign(value, { currency });

            return value;
          });
        });
      })
      .otherwise(() => undefined);

    store.ui.commandMenu.setOpen(false);
  };

  return (
    <>
      <CommandSubItem
        rightLabel='USD'
        icon={<CurrencyDollar />}
        leftLabel='Change currency'
        onSelectAction={() => handleSelect(Currency.Usd)}
        rightAccessory={
          opportunity?.value?.currency === Currency.Usd && <Check />
        }
      />
      <CommandSubItem
        rightLabel='EUR'
        icon={<CurrencyEuro />}
        leftLabel='Change currency'
        onSelectAction={() => handleSelect(Currency.Eur)}
        rightAccessory={
          opportunity?.value?.currency === Currency.Eur && <Check />
        }
      />
      <CommandSubItem
        rightLabel='GBP'
        icon={<CurrencyPound />}
        leftLabel='Change currency'
        onSelectAction={() => handleSelect(Currency.Gbp)}
        rightAccessory={
          opportunity?.value?.currency === Currency.Gbp && <Check />
        }
      />
    </>
  );
});

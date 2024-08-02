import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const UpdateHealthStatus = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = store.organizations.value.get(context.id as string);
  const label = `Organization - ${entity?.value?.name}`;

  const handleSelect =
    (renewalLikelihood: OpportunityRenewalLikelihood) => () => {
      if (!context.id) return;

      if (!entity) return;
      entity?.update((value) => {
        return {
          ...value,
          accountDetails: {
            ...value.accountDetails,
            renewalSummary: {
              ...value.accountDetails?.renewalSummary,
              renewalLikelihood,
            },
          },
        };
      });

      store.ui.commandMenu.toggle('RenameOrganizationProperty');
    };

  return (
    <Command label='Change health status...'>
      <CommandInput label={label} placeholder='Change health status...' />

      <Command.List>
        <CommandItem
          key={OpportunityRenewalLikelihood.HighRenewal}
          onSelect={handleSelect(OpportunityRenewalLikelihood.HighRenewal)}
          rightAccessory={
            entity?.value.accountDetails?.renewalSummary?.renewalLikelihood ===
            OpportunityRenewalLikelihood.HighRenewal ? (
              <Check />
            ) : null
          }
        >
          <span className='text-greenLight-500'>High</span>
        </CommandItem>
        <CommandItem
          key={OpportunityRenewalLikelihood.MediumRenewal}
          onSelect={handleSelect(OpportunityRenewalLikelihood.MediumRenewal)}
          rightAccessory={
            entity?.value.accountDetails?.renewalSummary?.renewalLikelihood ===
            OpportunityRenewalLikelihood.MediumRenewal ? (
              <Check />
            ) : null
          }
        >
          <span className='text-warning-500'>Medium</span>
        </CommandItem>
        <CommandItem
          key={OpportunityRenewalLikelihood.LowRenewal}
          onSelect={handleSelect(OpportunityRenewalLikelihood.LowRenewal)}
          rightAccessory={
            entity?.value.accountDetails?.renewalSummary?.renewalLikelihood ===
            OpportunityRenewalLikelihood.LowRenewal ? (
              <Check />
            ) : null
          }
        >
          <span className='text-error-500'>Low</span>
        </CommandItem>
      </Command.List>
    </Command>
  );
});

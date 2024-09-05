import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const UpdateHealthStatus = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = match(context.entity)
    .returnType<OrganizationStore | OrganizationStore[] | undefined>()

    .with('Organization', () =>
      store.organizations.value.get(context.ids?.[0] as string),
    )
    .with(
      'Organizations',
      () =>
        context.ids?.map((e: string) =>
          store.organizations.value.get(e),
        ) as OrganizationStore[],
    )
    .otherwise(() => undefined);

  const label = match(context.entity)
    .with(
      'Organization',
      () => `Organization - ${(entity as OrganizationStore)?.value?.name}`,
    )
    .with('Organizations', () => `${context.ids?.length} organizations`)

    .otherwise(() => '');

  const handleSelect =
    (renewalLikelihood: OpportunityRenewalLikelihood) => () => {
      if (!context.ids?.[0]) return;

      if (!entity) return;

      match(context.entity)
        .with('Organization', () => {
          (entity as OrganizationStore)?.update((value) => {
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
        })
        .with('Organizations', () =>
          store.organizations.updateHealth(
            context.ids as string[],
            renewalLikelihood,
          ),
        )
        .otherwise(() => undefined);

      store.ui.commandMenu.toggle('RenameOrganizationProperty');
    };

  const healthStatus =
    context.entity === 'Organization' &&
    (entity as OrganizationStore)?.value.accountDetails?.renewalSummary
      ?.renewalLikelihood;

  return (
    <Command
      label='Change health status...'
      onKeyDown={(e) => {
        e.stopPropagation();
      }}
    >
      <CommandInput
        label={label}
        placeholder='Change health status...'
        onKeyDownCapture={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            store.ui.commandMenu.setOpen(false);
          }
        }}
      />

      <Command.List>
        <CommandItem
          key={OpportunityRenewalLikelihood.HighRenewal}
          onSelect={handleSelect(OpportunityRenewalLikelihood.HighRenewal)}
          rightAccessory={
            healthStatus === OpportunityRenewalLikelihood.HighRenewal ? (
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
            healthStatus === OpportunityRenewalLikelihood.MediumRenewal ? (
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
            healthStatus === OpportunityRenewalLikelihood.LowRenewal ? (
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

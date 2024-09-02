import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import { OrganizationStore } from '@store/Organizations/Organization.store';
import { getStageFromColumn } from '@opportunities/components/ProspectsBoard/columns';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import {
  InternalStage,
  OrganizationStage,
  OrganizationRelationship,
} from '@graphql/types';
import {
  stageOptions,
  getStageOptions,
} from '@organization/components/Tabs/panels/AboutPanel/util';

type OpportunityStage = InternalStage | string;

export const ChangeStage = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const opportunityStages = store.tableViewDefs
    .getById(store.tableViewDefs.opportunitiesPreset ?? '')
    ?.value?.columns?.map((column) => ({
      value: getStageFromColumn(column),
      label: column.name,
    }));

  const entity = match(context.entity)
    .returnType<
      | OpportunityStore
      | OrganizationStore
      | OrganizationStore[]
      | OpportunityStore[]
      | undefined
    >()
    .with('Opportunity', () =>
      store.opportunities.value.get(context.ids?.[0] as string),
    )
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
    .with(
      'Opportunities',
      () =>
        context.ids?.map((e: string) =>
          store.opportunities.value.get(e),
        ) as OpportunityStore[],
    )
    .otherwise(() => undefined);

  const label = match(context.entity)
    .with(
      'Organization',
      () => `Organization - ${(entity as OrganizationStore)?.value?.name}`,
    )
    .with('Organizations', () => `${context.ids?.length} organizations`)
    .with(
      'Opportunity',
      () => `Opportunity - ${(entity as OpportunityStore)?.value?.name}`,
    )
    .with('Opportunities', () => `${context.ids?.length} opportunities`)
    .otherwise(() => '');

  const selectedStageOption = match(context.entity)
    .with('Organization', () =>
      stageOptions.find(
        (option) => option.value === (entity as OrganizationStore)?.value.stage,
      ),
    )
    .with('Opportunity', () =>
      opportunityStages?.find(
        (option) =>
          option.value === (entity as OpportunityStore)?.value?.externalStage ||
          option.value === (entity as OpportunityStore)?.value?.internalStage,
      ),
    )
    .with('Opportunities', () =>
      opportunityStages?.find(
        (option) =>
          option.value === (entity as OpportunityStore)?.value?.externalStage ||
          option.value === (entity as OpportunityStore)?.value?.internalStage,
      ),
    )
    .otherwise(() => undefined);

  const applicableStageOptions = match(context.entity)
    .with('Organization', () =>
      getStageOptions((entity as OrganizationStore).value?.relationship),
    )
    .with('Organizations', () =>
      getStageOptions(OrganizationRelationship.Prospect),
    )
    .with('Opportunity', () => opportunityStages ?? [])
    .with('Opportunities', () => opportunityStages ?? [])
    .otherwise(() => []);

  const handleSelect = (value: OrganizationStage | OpportunityStage) => () => {
    if (!context.ids?.[0]) return;

    if (!entity) return;

    match(context.entity)
      .with('Organization', () => {
        (entity as OrganizationStore)?.update((org) => {
          org.stage = value as OrganizationStage;

          return org;
        });
      })
      .with('Organizations', () => {
        store.organizations.updateStage(
          context.ids as string[],
          value as OrganizationStage,
        );
      })
      .with('Opportunities', () => {
        store.opportunities.updateStage(context.ids as string[], value);
      })
      .with('Opportunity', () => {
        (entity as OpportunityStore)?.update((opp) => {
          if (
            [InternalStage.ClosedLost, InternalStage.ClosedWon].includes(
              value as InternalStage,
            )
          ) {
            opp.internalStage = value as InternalStage;

            return opp;
          }
          opp.internalStage = InternalStage.Open;
          opp.externalStage = value;

          return opp;
        });
      });

    store.ui.commandMenu.toggle('ChangeStage');
  };

  return (
    <Command label='Change Stage'>
      <CommandInput label={label} placeholder='Change stage...' />

      <Command.List>
        {applicableStageOptions.map((option) => (
          <CommandItem
            key={option.value}
            onSelect={handleSelect(option.value)}
            rightAccessory={
              selectedStageOption?.value === option.value ? <Check /> : null
            }
          >
            {option.label}
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import { OrganizationStore } from '@store/Organizations/Organization.store';
import { getStageFromColumn } from '@opportunities/components/ProspectsBoard/columns';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { InternalStage, OrganizationStage } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import {
  stageOptions,
  getStageOptions,
} from '@organization/components/Tabs/panels/AboutPanel/util.ts';

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
    .returnType<OpportunityStore | OrganizationStore | undefined>()
    .with('Opportunity', () =>
      store.opportunities.value.get(context.id as string),
    )
    .with('Organization', () =>
      store.organizations.value.get(context.id as string),
    )
    .otherwise(() => undefined);

  const label = match(context.entity)
    .with('Organization', () => `Organization - ${entity?.value?.name}`)
    .with('Opportunity', () => `Opportunity - ${entity?.value?.name}`)
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
    .otherwise(() => undefined);

  const applicableStageOptions = match(context.entity)
    .with('Organization', () =>
      getStageOptions((entity as OrganizationStore).value?.relationship),
    )
    .with('Opportunity', () => opportunityStages ?? [])
    .otherwise(() => []);

  const handleSelect = (value: OrganizationStage | OpportunityStage) => () => {
    if (!context.id) return;

    if (!entity) return;

    match(context.entity)
      .with('Organization', () => {
        (entity as OrganizationStore)?.update((org) => {
          org.stage = value as OrganizationStage;

          return org;
        });
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

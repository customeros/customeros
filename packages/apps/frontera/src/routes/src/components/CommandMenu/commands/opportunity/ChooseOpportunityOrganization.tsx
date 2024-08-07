import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { InternalType, InternalStage } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChooseOpportunityOrganization = observer(() => {
  const store = useStore();
  const [search, setSearch] = useState('');
  const context = store.ui.commandMenu.context;

  const filteredOrganizations =
    search.length > 0
      ? store.organizations.toComputedArray((arr) => {
          return arr.filter((org) => org.value.name.includes(search));
        })
      : [];

  const handleSelect = (orgId: string) => () => {
    const organization = store.organizations.value.get(orgId)?.value;
    const stage = context?.meta?.stage;

    if (!organization || !stage) return;

    const isInternalStage =
      stage === InternalStage.ClosedLost || stage === InternalStage.ClosedWon;

    store.opportunities.create({
      organization,
      name: `${organization.name}'s opportunity`,
      internalType: InternalType.Nbo,
      externalStage: isInternalStage ? '' : stage,
      internalStage: isInternalStage ? stage : InternalStage.Open,
    });

    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('OpportunityHub');
  };

  return (
    <Command>
      <CommandInput
        value={search}
        label='Organization'
        placeholder='Choose organization'
        onValueChange={(v) =>
          setSearch(v.normalize('NFD').replace(/[\u0300-\u036f]/g, ''))
        }
      />

      <Command.List>
        {filteredOrganizations.map((org) => (
          <CommandItem key={org.getId()} onSelect={handleSelect(org.getId())}>
            {org.value.name}
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});

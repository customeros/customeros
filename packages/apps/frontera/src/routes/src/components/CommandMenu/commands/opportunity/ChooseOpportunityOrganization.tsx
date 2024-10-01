import { useMemo, useState } from 'react';

import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';

import { Avatar } from '@ui/media/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { InternalType, InternalStage } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

interface Organization {
  id: string;
  name: string;
  logo?: string;
}

export const ChooseOpportunityOrganization = observer(() => {
  const store = useStore();
  const [search, setSearch] = useState('');
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [filteredOrganizations, setFilteredOrganizations] = useState<any>([]);
  const context = store.ui.commandMenu.context;

  const organizationsList = useMemo(
    () =>
      store.organizations.toArray().map((org) => ({
        id: org.value.metadata.id,
        name: org.value.name,
        logo: org.value.logo,
      })),
    [store.organizations],
  );

  const fuse = useMemo(() => {
    return new Fuse(organizationsList, {
      keys: ['name'],
      threshold: 0.3,
    });
  }, [organizationsList]);

  const handleSearch = (v: string) => {
    const normalizedValue = v
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '');

    setSearch(normalizedValue);

    if (normalizedValue.length > 0) {
      const results = fuse.search(normalizedValue, { limit: 20 });

      setFilteredOrganizations(results.map((v) => v.item));
    } else {
      setFilteredOrganizations([]);
    }
  };

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
    <Command shouldFilter={false}>
      <CommandInput
        value={search}
        label='Organization'
        onValueChange={handleSearch}
        placeholder='Choose organization'
        dataTest='opp-kanban-choose-organization'
      />

      <Command.List>
        {filteredOrganizations.map((org: Organization) => (
          <CommandItem key={org.id} onSelect={handleSelect(org.id)}>
            <div className='flex items-center'>
              <Avatar
                size='xxs'
                name={org.name}
                className='mr-2'
                src={org.logo || ''}
                variant='outlineSquare'
              />
              {org.name}
            </div>
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});

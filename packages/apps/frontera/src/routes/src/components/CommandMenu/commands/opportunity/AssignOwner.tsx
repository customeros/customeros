import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { Check } from '@ui/media/icons/Check';
import { User01 } from '@ui/media/icons/User01';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandItem,
  CommandInput,
  useCommandState,
} from '@ui/overlay/CommandMenu';

export const AssignOwner = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const users = store.users.toArray();

  const entity = match(context.entity)
    .with('Opportunity', () =>
      store.opportunities.value.get(context.id as string),
    )
    .with('Organization', () =>
      store.organizations.value.get(context.id as string),
    )
    .otherwise(() => undefined);
  const label = match(context.entity)
    .with('Opportunity', () => `Opportunity - ${entity?.value?.name}`)
    .with('Organization', () => `Organization - ${entity?.value?.name}`)
    .otherwise(() => undefined);

  const handleSelect = (userId: string) => () => {
    if (!context.id) return;
    const user = store.users.value.get(userId);

    if (!user) return;

    match(context.entity)
      .with('Opportunity', () => {
        if (!entity) return;
        (entity as OpportunityStore)?.update((value) => {
          if (!value.owner) {
            Object.assign(value, { owner: user.value });

            return value;
          }

          Object.assign(value.owner, user.value);

          return value;
        });
      })
      .with('Organization', () => {
        if (!entity) return;
        (entity as OrganizationStore)?.update((value) => {
          if (!value.owner) {
            Object.assign(value, { owner: user.value });

            return value;
          }

          Object.assign(value.owner, user.value);

          return value;
        });
      })
      .otherwise(() => {});

    store.ui.commandMenu.toggle('AssignOwner');
  };

  return (
    <Command label='Pick Owner'>
      <CommandInput label={label} placeholder='Assign owner...' />

      <Command.List>
        <EmptySearch />

        {users.map((user) => (
          <CommandItem
            key={user.id}
            onSelect={handleSelect(user.id)}
            rightAccessory={user.id === entity?.owner?.id ? <Check /> : null}
            leftAccessory={
              <Avatar
                size='xs'
                textSize='xxs'
                name={user.name ?? 'Unnamed'}
                className='border border-gray-200'
                src={user.value.profilePhotoUrl ?? undefined}
                icon={<User01 className='text-gray-500 size-3' />}
              />
            }
          >
            {user.name}
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});

const EmptySearch = () => {
  const search = useCommandState((state) => state.search);

  return (
    <Command.Empty>{`No users found with name "${search}".`}</Command.Empty>
  );
};

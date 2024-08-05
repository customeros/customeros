import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { User01 } from '@ui/media/icons/User01.tsx';
import { Avatar } from '@ui/media/Avatar/Avatar.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const AssignOwner = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const users = store.users.toArray();

  const entity = match(context.entity)
    .returnType<
      OpportunityStore | OrganizationStore | OrganizationStore[] | undefined
    >()
    .with('Opportunity', () =>
      store.opportunities.value.get((context.ids as string[])?.[0]),
    )
    .with(
      'Organizations',
      () =>
        context.ids?.map((e: string) =>
          store.organizations.value.get(e),
        ) as OrganizationStore[],
    )
    .with('Organization', () =>
      store.organizations.value.get((context.ids as string[])?.[0]),
    )
    .otherwise(() => undefined);
  const label = match(context.entity)
    .with(
      'Opportunity',
      () => `Opportunity - ${(entity as OpportunityStore)?.value?.name}`,
    )
    .with(
      'Organization',
      () => `Organization - ${(entity as OrganizationStore)?.value?.name}`,
    )
    .with('Organizations', () => `${context.ids?.length} organizations`)
    .otherwise(() => undefined);

  const handleSelect = (userId: string) => () => {
    if (!context.ids?.[0]) return;
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
      .with('Organizations', () => {
        if (!(entity as OrganizationStore[])?.length) return;
        (entity as OrganizationStore[]).forEach((org) => {
          org.update((value) => {
            if (!value.owner) {
              Object.assign(value, { owner: user.value });

              return value;
            }

            Object.assign(value.owner, user.value);

            return value;
          });
        });
      })
      .otherwise(() => {});

    store.ui.commandMenu.toggle('AssignOwner');
  };

  return (
    <Command label='Pick Owner'>
      <CommandInput label={label} placeholder='Assign owner...' />

      <Command.List>
        {users.map((user) => (
          <CommandItem
            key={user.id}
            onSelect={handleSelect(user.id)}
            rightAccessory={
              user.id ===
              (entity as OrganizationStore | OpportunityStore)?.owner?.id ? (
                <Check />
              ) : null
            }
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

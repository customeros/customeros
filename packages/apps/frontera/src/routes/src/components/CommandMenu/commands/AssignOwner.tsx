import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

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
  const opportunity = store.opportunities.value.get(context.id as string);

  const label = match(context.entity)
    .with('Opportunity', () => `Opportunity - ${opportunity?.value?.name}`)
    .otherwise(() => undefined);

  const handleSelect = (userId: string) => () => {
    if (!context.id) return;
    const user = store.users.value.get(userId);

    if (!user) return;

    match(context.entity)
      .with('Opportunity', () => {
        if (!opportunity) return;
        opportunity?.update((value) => {
          if (!value.owner) {
            Object.assign(value, { owner: user.value });

            return value;
          }

          Object.assign(value.owner, user.value);

          return value;
        });
      })
      .with('Organization', () => {})
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
            leftAccessory={
              <Avatar
                size='xs'
                textSize='xxs'
                name={user.name ?? 'Unnamed'}
                src={user.value.profilePhotoUrl ?? undefined}
                icon={<User01 className='text-gray-500 size-3' />}
                className='border border-gray-200'
              />
            }
            rightAccessory={
              user.id === opportunity?.owner?.id ? <Check /> : null
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

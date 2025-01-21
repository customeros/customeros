import { observer } from 'mobx-react-lite';

import { User } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ConnectedUsersProps {
  users: User[];
}

export const ConnectedUsers = observer(({ users }: ConnectedUsersProps) => {
  const store = useStore();

  if (!users.length) return <p className='text-gray-400'>No one</p>;

  const user = users
    ?.map(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (l: any) => store.users.value.get(l.id)?.name,
    )
    .join(', ');

  return (
    <Tooltip
      label={
        users.length > 1
          ? users
              .slice(1, users.length)
              ?.map((e) => e?.name)
              .join(', ')
          : ''
      }
    >
      <div className='flex w-fit'>
        <div className='bg-gray-100 rounded-md w-fit px-1.5 '>{user}</div>
        {users?.length > 1 && (
          <div className='rounded-md w-fit px-1.5 ml-1 text-gray-500'>
            +{users?.length - 1}
          </div>
        )}
      </div>
    </Tooltip>
  );
});

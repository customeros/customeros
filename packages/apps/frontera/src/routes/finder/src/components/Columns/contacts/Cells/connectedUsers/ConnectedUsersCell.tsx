import { observer } from 'mobx-react-lite';

import { User } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';

interface ConnectedUsersProps {
  users: User[];
}

export const ConnectedUsers = observer(({ users }: ConnectedUsersProps) => {
  const store = useStore();

  if (!users.length) return <p className='text-gray-400'> No one </p>;

  const usersDisplayed = users?.map(
    (l: User) => store.users.value.get(l.id)?.name,
  );

  return (
    <p
      className='flex w-fit gap-x-1'
      title={
        usersDisplayed?.length > 1
          ? usersDisplayed?.map((e) => e).join(', ')
          : ''
      }
    >
      {usersDisplayed?.map((name, i) => (
        <div
          key={`connected-user-${i}`}
          className='bg-gray-100 rounded-md w-fit px-1.5 '
        >
          {name}
        </div>
      ))}
    </p>
  );
});

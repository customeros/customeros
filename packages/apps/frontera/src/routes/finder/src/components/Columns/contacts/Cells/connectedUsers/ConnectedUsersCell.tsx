import { observer } from 'mobx-react-lite';

import { User } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ConnectedUsersProps {
  users: User[];
}

export const ConnectedUsers = observer(({ users }: ConnectedUsersProps) => {
  if (!users.length) return <p className='text-gray-400'>No one</p>;

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
        <div className='bg-gray-100 rounded-md w-fit px-1.5 '>
          {users?.[0]?.name}
        </div>
        {users?.length > 1 && (
          <div className='rounded-md w-fit px-1.5 ml-1 text-gray-500'>
            +{users?.length - 1}
          </div>
        )}
      </div>
    </Tooltip>
  );
});

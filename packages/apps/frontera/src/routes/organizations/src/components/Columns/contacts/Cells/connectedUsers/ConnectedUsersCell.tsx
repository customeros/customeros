import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ConnectedUsersProps {
  users: string[];
}

export const ConnectedUsers = observer(({ users }: ConnectedUsersProps) => {
  const store = useStore();
  const data = users.map((userId) => store.users.value.get(userId));

  if (!data.length) return <p className='text-gray-400'>No one</p>;

  return (
    <Tooltip
      label={
        data.length > 1
          ? data
              .slice(1, data.length)
              ?.map((e) => e?.name)
              .join(', ')
          : ''
      }
    >
      <div className='flex w-fit'>
        <div className='bg-gray-100 rounded-md w-fit px-1.5 '>
          {data?.[0]?.name}
        </div>
        {data?.length > 1 && (
          <div className='rounded-md w-fit px-1.5 ml-1 text-gray-500'>
            +{data?.length - 1}
          </div>
        )}
      </div>
    </Tooltip>
  );
});

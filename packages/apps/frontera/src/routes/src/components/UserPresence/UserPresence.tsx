import { useChannel } from '@shared/hooks/useChannel';

import { UserHexagon } from '../UserHexagon';

interface UserPresenceProps {
  channelName: string;
}

export const UserPresence = ({ channelName }: UserPresenceProps) => {
  const { presentUsers, currentUserId } = useChannel(channelName);

  return (
    <div className='flex'>
      {presentUsers.map(([user]) => (
        <UserHexagon
          id={user?.user_id}
          key={user?.user_id}
          color={user?.color}
          name={user?.username}
          isCurrent={user?.user_id === currentUserId}
        />
      ))}
    </div>
  );
};

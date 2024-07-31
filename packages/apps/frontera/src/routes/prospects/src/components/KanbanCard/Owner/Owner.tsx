import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { User01 } from '@ui/media/icons/User01';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface OwnerProps {
  ownerId?: string | null;
}

export const Owner = observer(({ ownerId }: OwnerProps) => {
  const store = useStore();
  const user = store.users.value.get(ownerId ?? '');

  return (
    <Tooltip label={user?.name}>
      <Avatar
        size='xs'
        textSize='xxs'
        name={user?.name ?? 'Unnamed'}
        src={user?.value?.profilePhotoUrl ?? ''}
        icon={<User01 className='text-gray-500 size-3' />}
        className={cn(
          'w-5 h-5 min-w-5',
          user?.value?.profilePhotoUrl ? '' : 'border border-gray-200',
        )}
      />
    </Tooltip>
  );
});

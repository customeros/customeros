import { observer } from 'mobx-react-lite';

import { Avatar } from '@ui/media/Avatar';
import { User01 } from '@ui/media/icons/User01';
import { useStore } from '@shared/hooks/useStore';

interface MailboxCellProps {
  mailbox: string;
}

export const MailboxCell = observer(({ mailbox }: MailboxCellProps) => {
  const store = useStore();
  const mailboxStore = store.mailboxes.value.get(mailbox);
  const userId = mailboxStore?.value.userId;

  const user = userId ? store.users.value.get(userId) : null;

  return (
    <div className='flex gap-1 items-center'>
      <Avatar
        size='xs'
        textSize='xxs'
        variant='circle'
        name={user?.name ?? ''}
        src={user?.value?.profilePhotoUrl ?? ''}
        icon={<User01 className='text-gray-500 size-3' />}
        className={'w-5 h-5 min-w-5 mr-2 border border-gray-200'}
      />
      <p>{mailbox}</p>
    </div>
  );
});

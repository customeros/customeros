import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';

interface OwnerProps {
  ownerId?: string;
}

export const OwnerCell = observer(({ ownerId }: OwnerProps) => {
  const store = useStore();
  const userOwner = store.users.value.get(ownerId ?? '');

  return (
    <div
      className='flex w-full gap-1 items-center cursor-pointer'
      onClick={() => {
        store.ui.commandMenu.setType('AssignOwner');
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <p
        data-test='organization-owner-in-all-orgs-table'
        className={cn(ownerId ? 'text-gray-700' : 'text-gray-400')}
      >
        {userOwner?.name ?? 'No owner'}
      </p>
    </div>
  );
});

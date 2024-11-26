import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { User01 } from '@ui/media/icons/User01.tsx';

interface AvatarCellProps {
  id: string;
  name: string;
  icon?: string | null;
  logo?: string | null;
  canNavigate?: boolean;
}

export const AvatarCell = observer(
  ({ name, id, icon, logo, canNavigate }: AvatarCellProps) => {
    const store = useStore();
    const src = icon || logo;
    const fullName = name || 'Unnamed';
    const contactStore = store.contacts.value.get(id);

    const isEnriching = contactStore?.isEnriching;

    return (
      <div className='items-center ml-[1px]'>
        <Avatar
          size='xs'
          textSize='xs'
          tabIndex={-1}
          icon={<User01 />}
          src={src || undefined}
          variant='outlineCircle'
          name={isEnriching ? '' : fullName}
          className={cn(
            'text-gray-700 cursor-pointer focus:outline-none',
            !canNavigate && 'cursor-default',
            isEnriching && 'animate-pulse',
          )}
          onClick={() => {
            if (
              store.ui.contactPreviewCardOpen === true &&
              store.ui.focusRow === id
            ) {
              store.ui.setContactPreviewCardOpen(false);
            } else {
              store.ui.setFocusRow(id);
              store.ui.setContactPreviewCardOpen(true);
            }
          }}
        />
      </div>
    );
  },
);

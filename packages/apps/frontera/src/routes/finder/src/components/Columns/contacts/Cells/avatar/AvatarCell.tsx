import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { useStore } from '@shared/hooks/useStore';

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

    return (
      <div className='items-center ml-[1px]'>
        <Avatar
          size='xs'
          textSize='xs'
          tabIndex={-1}
          name={fullName}
          src={src || undefined}
          variant='outlineCircle'
          className={cn(
            'text-gray-700 cursor-pointer focus:outline-none',
            !canNavigate && 'cursor-default',
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

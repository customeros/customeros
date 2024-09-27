import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit01 } from '@ui/media/icons/Edit01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

import { TagsCell } from '../../../shared/Cells';

interface ContactCardProps {
  id: string;
}

export const ContactsTagsCell = observer(({ id }: ContactCardProps) => {
  const store = useStore();
  const [isHovered, setIsHovered] = useState(false);
  const contactStore = store.contacts.value.get(id);

  const tags = (contactStore?.value?.tags ?? []).filter((d) => !!d?.name);

  return (
    <div
      className='flex items-center '
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onDoubleClick={() => {
        store.ui.commandMenu.setType('EditPersonaTag');
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <TagsCell tags={tags ?? []} />
      {isHovered && (
        <IconButton
          size='xxs'
          variant='ghost'
          className='ml-3'
          aria-label='Edit tags'
          icon={<Edit01 className='text-gray-500' />}
          onClick={() => {
            store.ui.commandMenu.setType('EditPersonaTag');
            store.ui.commandMenu.setOpen(true);
          }}
        />
      )}
    </div>
  );
});

import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit01 } from '@ui/media/icons/Edit01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

import { TagsCell } from '../../../shared/Cells';

interface OrgCardProps {
  id: string;
}

export const OrganizationsTagsCell = observer(({ id }: OrgCardProps) => {
  const store = useStore();
  const organizationstore = store.organizations.value.get(id);
  const [isHovered, setIsHovered] = useState(false);

  const tags = (organizationstore?.value?.tags ?? []).filter((d) => !!d?.name);

  return (
    <div
      className='flex items-center '
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onDoubleClick={() => {
        store.ui.commandMenu.setType('ChangeTags');
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <TagsCell tags={tags ?? []} isHovered={isHovered} />
      {isHovered && (
        <IconButton
          size='xxs'
          className=' '
          variant='ghost'
          aria-label='Edit tags'
          icon={<Edit01 className='text-gray-500' />}
          onClick={() => {
            store.ui.commandMenu.setType('ChangeTags');
            store.ui.commandMenu.setOpen(true);
          }}
        />
      )}
    </div>
  );
});

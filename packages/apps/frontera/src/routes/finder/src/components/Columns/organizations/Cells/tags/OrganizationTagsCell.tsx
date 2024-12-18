import { observer } from 'mobx-react-lite';

import type { Tag } from '@graphql/types';

import { useStore } from '@shared/hooks/useStore';

import { TagsCell } from '../../../shared/Cells';

interface OrganizationTagsCellProps {
  id: string;
}

export const OrganizationsTagsCell = observer(
  ({ id }: OrganizationTagsCellProps) => {
    const store = useStore();
    const entity = store.organizations.getById(id);

    const tags = (entity?.value?.tags ?? []).filter((d) => !!d?.name);

    return (
      <div
        className='flex items-center cursor-pointer'
        onClick={() => {
          store.ui.commandMenu.setType('ChangeTags');
          store.ui.commandMenu.setOpen(true);
        }}
      >
        <TagsCell tags={(tags ?? []) as Tag[]} />
      </div>
    );
  },
);

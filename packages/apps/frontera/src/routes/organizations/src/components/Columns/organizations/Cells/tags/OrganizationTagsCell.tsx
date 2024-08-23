import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { TagsCell } from '../../../shared/Cells';

interface OrgCardProps {
  id: string;
}

export const OrganizationsTagsCell = observer(({ id }: OrgCardProps) => {
  const store = useStore();
  const organizationstore = store.organizations.value.get(id);
  const tags = (organizationstore?.value?.tags ?? []).filter((d) => !!d?.name);

  return (
    <div
      className='cursor-pointer'
      onDoubleClick={() => {
        store.ui.commandMenu.setType('ChangeTags');
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <TagsCell tags={tags ?? []} />
    </div>
  );
});

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Tag } from '@shared/types/__generated__/graphql.types';

import { TagsCell } from '../../../shared/Cells';

interface ContactCardProps {
  id: string;
}

export const ContactsTagsCell = observer(({ id }: ContactCardProps) => {
  const store = useStore();
  const contactStore = store.contacts.value.get(id);

  const tags = (contactStore?.value?.tags ?? []).filter((d) => !!d?.name);

  return (
    <div
      className='flex items-center cursor-pointer'
      onClick={() => {
        store.ui.commandMenu.setType('EditPersonaTag');
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <TagsCell tags={(tags ?? []) as Tag[]} />
    </div>
  );
});

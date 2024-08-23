import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

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
      className='cursor-pointer'
      onDoubleClick={() => {
        store.ui.commandMenu.setType('EditPersonaTag');
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <TagsCell tags={tags ?? []} />
    </div>
  );
});

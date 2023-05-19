import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '@spaces/atoms/checkbox';
import { Contact } from '../../../graphQL/__generated__/generated';
import { getContactDisplayName } from '../../../utils';
import { LinkCell, TableCell } from '@spaces/atoms/table/table-cells/TableCell';
import { ContactAvatar } from '@spaces/molecules/contact-avatar/ContactAvatar';

export const ContactTableCell: React.FC<{
  contact?: Contact;
}> = ({ contact }) => {
  const mode = useRecoilValue(tableMode);
  const [selectedIds, setSelectedIds] = useRecoilState(selectedItemsIds);

  if (!contact) {
    return <span className={styles.emptyCell}>-</span>;
  }
  const handleCheckboxToggle = () => {
    const isChecked = selectedIds.findIndex((id) => contact.id === id) !== -1;

    if (isChecked) {
      const filteredIds = selectedIds.filter((id) => id !== contact.id);
      setSelectedIds([...filteredIds]);
      return;
    }

    setSelectedIds((oldSelectedIds) => {
      return Array.from(new Set([...oldSelectedIds, contact.id]));
    });
  };
  const displayName = getContactDisplayName(contact);

  return (
    <>
      {mode === 'MERGE' && (
        <Checkbox
          checked={selectedIds.findIndex((id) => contact.id === id) !== -1}
          type='checkbox'
          label={
            <TableCell label={displayName || 'Unnamed'}>
              <ContactAvatar contactId={contact.id} name={displayName} />
            </TableCell>
          }
          //@ts-expect-error fixme
          onChange={handleCheckboxToggle}
        />
      )}

      {mode !== 'MERGE' && (
        <LinkCell label={displayName} url={`/contact/${contact.id}`}>
          <ContactAvatar contactId={contact.id} name={displayName} />
        </LinkCell>
      )}
    </>
  );
};

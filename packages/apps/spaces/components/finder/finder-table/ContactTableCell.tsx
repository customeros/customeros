import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '@spaces/atoms/checkbox';
import { Contact } from '../../../graphQL/__generated__/generated';
import { getContactDisplayName } from '../../../utils';
import { LinkCell, TableCell } from '@spaces/atoms/table/table-cells/TableCell';

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

  return (
    <>
      {mode === 'MERGE_CONTACT' && (
        <div className={styles.mergableCell}>
          <Checkbox
            checked={selectedIds.findIndex((id) => contact.id === id) !== -1}
            type='checkbox'
            label={<TableCell label={getContactDisplayName(contact)} />}
            //@ts-expect-error fixme
            onChange={handleCheckboxToggle}
          />
        </div>
      )}

      {mode !== 'MERGE_CONTACT' && (
        <LinkCell
          label={getContactDisplayName(contact)}
          url={`/contact/${contact.id}`}
        />
      )}
    </>
  );
};

import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '../../ui-kit/atoms/input';
import { FinderCell } from './FinderTableCell';
import { Contact } from '../../../graphQL/__generated__/generated';
import { getContactDisplayName } from '../../../utils';

export const ContactTableCell: React.FC<{
  contact: Contact;
}> = ({ contact }) => {
  const mode = useRecoilValue(tableMode);
  const [selectedIds, setSelectedIds] = useRecoilState(selectedItemsIds);
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

  if (!contact) {
    return <span>-</span>;
  }

  return (
    <div className={styles.mergableCell}>
      <FinderCell
        label={getContactDisplayName(contact)}
        subLabel={
          (contact.jobRoles.find((role) => role.primary) || contact.jobRoles[0])
            ?.jobTitle || ''
        }
        url={`/contact/${contact.id}`}
      />
      <div className={styles.checkboxContainer}>
        {mode === 'MERGE_CONTACT' && (
          <Checkbox
            checked={selectedIds.findIndex((id) => contact.id === id) !== -1}
            onChange={() => handleCheckboxToggle()}
          />
        )}
      </div>
    </div>
  );
};

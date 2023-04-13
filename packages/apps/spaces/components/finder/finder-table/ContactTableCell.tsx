import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '../../ui-kit/atoms/input';
import { FinderCell } from './FinderTableCell';
import {
  Contact,
  Organization,
} from '../../../graphQL/__generated__/generated';
import { getContactDisplayName } from '../../../utils';

export const ContactTableCell: React.FC<{
  contact?: Contact;
  organization?: Organization;
}> = ({ contact, organization }) => {
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
    <div className={styles.mergableCell}>
      <div className={styles.checkboxContainer}>
        {mode === 'MERGE_CONTACT' && (
          <Checkbox
            checked={selectedIds.findIndex((id) => contact.id === id) !== -1}
            onChange={() => handleCheckboxToggle()}
          />
        )}
      </div>
      <div className={styles.finderCell}>
        <FinderCell
          label={getContactDisplayName(contact)}
          subLabel={
            (
              contact.jobRoles.find((role) =>
                organization?.id
                  ? role.organization?.id === organization.id
                  : role.primary,
              ) || contact.jobRoles[0]
            )?.jobTitle || ''
          }
          url={`/contact/${contact.id}`}
        />
      </div>
    </div>
  );
};

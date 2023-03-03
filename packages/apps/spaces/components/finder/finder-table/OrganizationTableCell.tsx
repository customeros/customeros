import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '../../ui-kit/atoms/input';
import { FinderCell } from './FinderTableCell';
import { Organization } from '../../../graphQL/__generated__/generated';

//@ts-expect-error fixme later
export const OrganizationTableCell: React.FC<{
  organization: Organization;
}> = ({ organization }) => {
  const mode = useRecoilValue(tableMode);
  const [selectedIds, setSelectedIds] = useRecoilState(selectedItemsIds);
  const handleCheckboxToggle = () => {
    const isChecked =
      selectedIds.findIndex((id) => organization.id === id) !== -1;

    if (isChecked) {
      const filteredIds = selectedIds.filter((id) => id !== organization.id);
      setSelectedIds([...filteredIds]);
      return;
    }

    setSelectedIds((oldSelectedIds) => {
      return Array.from(new Set([...oldSelectedIds, organization.id]));
    });
  };

  if (organization) {
    const industry = (
      <span className={'capitalise'}>
        {(organization?.industry ?? '')?.split('_').join(' ').toLowerCase()}
      </span>
    );
    return (
      <div className={styles.mergableCell}>
        <FinderCell
          label={organization.name}
          subLabel={industry}
          url={`/organization/${organization.id}`}
        />
        <div className={styles.checkboxContainer}>
          {mode === 'MERGE_ORG' && (
            <Checkbox
              checked={
                selectedIds.findIndex((id) => organization.id === id) !== -1
              }
              onChange={() => handleCheckboxToggle()}
            />
          )}
        </div>
      </div>
    );
  }
  if (!organization) {
    return <span>-</span>;
  }
};

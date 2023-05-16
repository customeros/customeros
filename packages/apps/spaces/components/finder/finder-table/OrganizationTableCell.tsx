import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '@spaces/atoms/checkbox';
import { Organization } from '../../../graphQL/__generated__/generated';
import { TableCell } from '@spaces/atoms/table';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { FinderCell } from '@spaces/finder/finder-table/FinderTableCell';

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
  if (!organization) {
    return <div className={styles.emptyCell}>-</div>;
  }

  const industry = (
    <span className={'capitalise'}>
      {(organization?.industry ?? '')?.split('_').join(' ').toLowerCase()}
    </span>
  );
  return (
    <>
      {mode === 'MERGE_ORG' && (
        <div className={styles.mergableCell}>
          <Checkbox
            type='checkbox'
            checked={
              selectedIds.findIndex((id) => organization.id === id) !== -1
            }
            label={
              <TableCell
                label={organization.name || 'Unnamed'}
                subLabel={organization.industry}
              />
            }
            //@ts-expect-error fixme
            onChange={() => handleCheckboxToggle()}
          />
        </div>
      )}

      {mode !== 'MERGE_ORG' && (
        <LinkCell
          label={
            organization.name && organization.name !== ''
              ? organization.name
              : 'Unnamed'
          }
          subLabel={industry}
          url={`/organization/${organization.id}`}
        />
      )}

      <div className={styles.finderCell}>
        <FinderCell
          label={
            organization.name && organization.name !== ''
              ? organization.name
              : 'Unnamed'
          }
          subLabel={industry}
          url={`/organization/${organization.id}`}
        />
      </div>
    </>
  );
};

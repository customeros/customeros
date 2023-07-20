import React from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Checkbox } from '@spaces/atoms/checkbox';
import { Organization } from '../../../graphQL/__generated__/generated';
import { TableCell } from '@spaces/atoms/table';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { OrganizationAvatar } from '@spaces/molecules/organization-avatar/OrganizationAvatar';

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

  const hasParent = !!organization.subsidiaryOf?.length;

  return (
    <>
      {mode === 'MERGE' && (
        <Checkbox
          type='checkbox'
          checked={selectedIds.findIndex((id) => organization.id === id) !== -1}
          label={
            <TableCell
              label={
                hasParent
                  ? organization.subsidiaryOf[0].organization.name
                  : organization.name || 'Unnamed'
              }
              subLabel={hasParent ? organization.name || 'Unnamed' : ''}
            >
              <OrganizationAvatar name={organization?.name || 'Unnamed'} />
            </TableCell>
          }
          //@ts-expect-error fixme
          onChange={() => handleCheckboxToggle()}
        />
      )}

      {mode !== 'MERGE' && (
        <LinkCell
          label={
            hasParent
              ? organization.subsidiaryOf[0].organization.name
              : organization.name || 'Unnamed'
          }
          subLabel={hasParent ? organization.name || 'Unnamed' : ''}
          url={`/organizations/${organization.id}?tab=about`}
        >
          <OrganizationAvatar name={organization?.name} />
        </LinkCell>
      )}
    </>
  );
};

import React, { useRef } from 'react';
import { DashboardTableAddressCell } from '@spaces/atoms/table/table-cells/TableCell';
import { Button } from '@spaces/atoms/button';
import { OverlayPanel } from '@spaces/atoms/overlay-panel';
import styles from './finder-table.module.scss';
import { FinderCell } from '@spaces/finder/finder-table/FinderTableCell';

export const OwnerTableCell = ({ owner }: { owner: any }) => {
  const op = useRef(null);

  return (
    <FinderCell label={owner ? owner.firstName + ' ' + owner.lastName : ''} />
  );
};

import React, { ReactNode } from 'react';
import { TableCell } from '@spaces/atoms/table';

export const FinderCell = ({
  label,
  subLabel,
}: {
  label: string | ReactNode;
  subLabel?: string | ReactNode;
}) => {
  return <TableCell label={label} subLabel={subLabel} />;
};

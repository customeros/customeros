import React, { ReactNode } from 'react';
import { TableCell } from '@spaces/atoms/table';

export const FinderCell = ({
  label,
  subLabel,
}: {
  label: string;
  subLabel?: string | ReactNode;
}) => {
  return <TableCell label={label} subLabel={subLabel} />;
};

import React, { ReactNode } from 'react';
import { TableCell } from '@spaces/atoms/table';

export const FinderCell = ({
  label,
  subLabel,
  url,
}: {
  label: string;
  subLabel?: string | ReactNode;
  url?: string;
}) => {
  return (
    <TableCell
      label={label}
      subLabel={subLabel}
      url={url}
    />
  );
};

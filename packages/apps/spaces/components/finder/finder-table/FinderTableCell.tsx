import React, { ReactNode } from 'react';
import { TableCell } from '@spaces/atoms/table';
import { Highlight } from '@spaces/atoms/highlight';

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
      label={<Highlight text={label}
                        highlight={''} />}
      subLabel={subLabel}
      url={url}
    />
  );
};

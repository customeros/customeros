import React, { ReactNode } from 'react';
import { useRecoilValue } from 'recoil';
import { finderSearchTerm } from '../../../state';
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
  const searchTern = useRecoilValue(finderSearchTerm);

  return (
    <TableCell
      label={<Highlight text={label} highlight={searchTern} />}
      subLabel={subLabel}
      url={url}
    />
  );
};

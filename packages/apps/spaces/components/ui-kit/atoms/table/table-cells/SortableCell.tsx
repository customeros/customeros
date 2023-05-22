import React from 'react';
import { IconButton } from '@spaces/atoms/icon-button';
import { SortingDirection } from '../../../../../graphQL/__generated__/generated';
import Sort from '@spaces/atoms/icons/Sort';

interface SortableCellProps {
  sort: any;
  column: any;
  setSortingState: any;
}

export const SortableCell: React.FC<SortableCellProps> = ({
  sort,
  setSortingState,
  column,
}) => {
  return (
    <IconButton
      isSquare
      mode='text'
      onClick={() => {
        setSortingState({
          direction:
            sort.direction === SortingDirection.Asc
              ? SortingDirection.Desc
              : SortingDirection.Asc,
          column,
        });
      }}
      label='Sort'
      size={'xxxxs'}
      icon={
        <Sort
          height={10}
          color='#969696'
          style={{
            transform:
              sort.column === column && sort.direction === SortingDirection.Asc
                ? 'rotate(180deg)'
                : '',
          }}
        />
      }
    />
  );
};

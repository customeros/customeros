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
      onDoubleClick={() =>
        setSortingState({
          direction: SortingDirection.Asc,
          column: undefined,
        })
      }
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
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <Sort
            height={8}
            color={
              sort.column === column && sort.direction !== SortingDirection.Asc
                ? '#3a3a3a'
                : '#969696'
            }
            style={{
              transform: 'rotate(180deg)',
              marginBottom: 2,
            }}
          />
          <Sort
            height={8}
            color={
              sort.column === column && sort.direction !== SortingDirection.Desc
                ? '#3a3a3a'
                : '#969696'
            }
          />
        </div>
      }
    />
  );
};

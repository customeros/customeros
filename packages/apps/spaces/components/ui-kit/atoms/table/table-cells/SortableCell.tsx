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
  const handleSortModeChange = () => {
    if (sort.direction === undefined) {
      setSortingState({
        direction: SortingDirection.Asc,
        column,
      });
      return;
    }
    if (sort.direction === SortingDirection.Asc) {
      setSortingState({
        direction: SortingDirection.Desc,
        column,
      });
      return;
    }
    if (sort.direction === SortingDirection.Desc) {
      setSortingState({
        direction: undefined,
        column: undefined,
      });
      return;
    }
  };
  return (
    <IconButton
      isSquare
      mode='text'
      onClick={handleSortModeChange}
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

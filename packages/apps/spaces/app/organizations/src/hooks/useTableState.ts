import { useQuery } from '@tanstack/react-query';
import { ColumnFiltersState, SortingState } from '@tanstack/react-table';

type TableState = {
  columnFilters: ColumnFiltersState;
  columnSorting: SortingState;
};

const fetcher = (): Promise<TableState> => {
  return new Promise((resolve) => {
    resolve({
      columnFilters: [],
      columnSorting: [],
    });
  });
};

export const useTableState = () => {
  return useQuery<TableState, unknown, TableState>({
    queryKey: ['tableState'],
    queryFn: fetcher,
  });
};

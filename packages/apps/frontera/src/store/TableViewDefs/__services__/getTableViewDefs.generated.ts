import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type TableViewDefsQueryVariables = Types.Exact<{ [key: string]: never }>;

export type TableViewDefsQuery = {
  __typename?: 'Query';
  tableViewDefs: Array<{
    __typename?: 'TableViewDef';
    id: string;
    name: string;
    tableType: Types.TableViewType;
    tableId: Types.TableIdType;
    order: number;
    icon: string;
    filters: string;
    sorting: string;
    isPreset: boolean;
    isShared: boolean;
    createdAt: any;
    defaultFilters: string;
    updatedAt: any;
    columns: Array<{
      __typename?: 'ColumnView';
      columnId: number;
      columnType: Types.ColumnViewType;
      name: string;
      width: number;
      visible: boolean;
      filter: string;
    }>;
  }>;
};

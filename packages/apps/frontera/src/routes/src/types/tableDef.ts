import { Node, Scalars, ColumnViewType } from '@graphql/types';

export type TableViewDef = Node & {
  tableType: TableViewType;
  columns: Array<ColumnView>;
  __typename?: 'TableViewDef';
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  icon: Scalars['String']['output'];
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  filters: Scalars['String']['output'];
  sorting: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type ColumnView = {
  __typename?: 'ColumnView';
  columnType: ColumnViewType;
  width: Scalars['Int']['output'];
};

export enum TableViewType {
  Invoices = 'INVOICES',
  Organizations = 'ORGANIZATIONS',
  Renewals = 'RENEWALS',
}

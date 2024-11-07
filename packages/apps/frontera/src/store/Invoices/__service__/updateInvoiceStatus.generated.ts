import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type UpdateInvoiceStatusMutationVariables = Types.Exact<{
  input: Types.InvoiceUpdateInput;
}>;

export type UpdateInvoiceStatusMutation = {
  __typename?: 'Mutation';
  invoice_Update: {
    __typename?: 'Invoice';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};

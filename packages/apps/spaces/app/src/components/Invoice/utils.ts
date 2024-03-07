import { ComparisonOperator } from '@graphql/types';

export const filterOutDryRunInvoices = {
  filter: {
    property: 'DRY_RUN',
    operation: ComparisonOperator.Eq,
    value: false,
  },
};

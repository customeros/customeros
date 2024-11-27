import { ComparisonOperator } from '@graphql/types';

export const filterOutDryRunInvoices = {
  AND: [
    {
      filter: {
        property: 'DRY_RUN',
        operation: ComparisonOperator.Equals,
        value: false,
      },
    },
  ],
};

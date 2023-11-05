import { Filter, ComparisonOperator } from '@graphql/types';

export const mapGCliSearchTermsToFilterList = (
  searchTerms: any[],
  searchFor: string,
): Filter[] => {
  const filters = [] as Filter[];
  searchTerms.forEach((item: any) => {
    if (item.type === 'STATE') {
      filters.push({
        filter: {
          property: 'REGION',
          operation: ComparisonOperator.Eq,
          value: item.display,
        },
      });
      filters.push({
        filter: {
          property: 'REGION',
          operation: ComparisonOperator.Eq,
          value: item.data[0].value,
        },
      });
    } else {
      filters.push({
        filter: {
          property: searchFor,
          operation: ComparisonOperator.Eq,
          value: item.display,
        },
      });
    }
  });

  return filters;
};

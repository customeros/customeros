// import { produce } from 'immer';

// import { getServerGraphQLClient } from '@shared/util/getServerGraphQLClient';
// import {
//   Filter,
//   SortBy,
//   SortingDirection,
//   ComparisonOperator,
// } from '@graphql/types';
// import {
//   GlobalCacheQuery,
//   GlobalCacheDocument,
//   GlobalCacheQueryVariables,
// } from '@shared/graphql/global_Cache.generated';
// import {
//   GetOrganizationsQuery,
//   GetOrganizationsDocument,
//   GetOrganizationsQueryVariables,
// } from '@organizations/graphql/getOrganizations.generated';

import { Search } from './src/components/Search';
import { OrganizationsTable } from './src/components/OrganizationsTable';

export const OrganizationsPage = () => {
  // const client = getServerGraphQLClient();
  // const { preset, searchTerm } = searchParams;

  // let initialData: GetOrganizationsQuery | undefined = undefined;

  // try {
  //   const globalCache = await client.request<
  //     GlobalCacheQuery,
  //     GlobalCacheQueryVariables
  //   >(GlobalCacheDocument);
  //   const userId = globalCache?.global_Cache?.user?.id;

  //   const filters = createFilters({ preset, userId, searchTerm });

  //   initialData = await client.request<
  //     GetOrganizationsQuery,
  //     GetOrganizationsQueryVariables
  //   >(GetOrganizationsDocument, {
  //     pagination: { limit: 40, page: 1 },
  //     ...filters,
  //   });
  // } catch (e) {
  //   console.error('Failed to fetch initial Organizations data', e);
  // }

  return (
    <>
      <Search />
      <OrganizationsTable />
    </>
  );
};

// function createFilters({
//   preset,
//   userId,
//   searchTerm,
// }: {
//   userId: string;
//   preset?: string;
//   searchTerm?: string;
// }) {
//   const sort: SortBy = {
//     by: 'LAST_TOUCHPOINT',
//     direction: SortingDirection.Desc,
//     caseSensitive: false,
//   };

//   const where = (() => {
//     return produce<Filter>({ AND: [] }, (draft) => {
//       if (!draft.AND) {
//         draft.AND = [];
//       }

//       if (searchTerm) {
//         draft.AND.push({
//           filter: {
//             property: 'ORGANIZATION',
//             value: searchTerm,
//             operation: ComparisonOperator.Contains,
//             caseSensitive: false,
//           },
//         });
//       }

//       if (preset) {
//         const [property, value] = (() => {
//           if (preset === 'customer') {
//             return ['IS_CUSTOMER', [true]];
//           }
//           if (preset === 'portfolio') {
//             return ['OWNER_ID', [userId]];
//           }

//           return [];
//         })();
//         if (!property || !value) return;
//         draft.AND.push({
//           filter: {
//             property,
//             value,
//             operation: ComparisonOperator.Eq,
//             includeEmpty: false,
//           },
//         });
//       }
//     });
//   })();

//   return {
//     where,
//     sort,
//   };
// }

import React, { useEffect, useLayoutEffect, useState } from 'react';
import styles from './organization-list.module.scss';
import { organizationListColumns } from './columns/OrganizationListColumns';
import { useFinderOrganizationTableData } from '@spaces/hooks/useFinderOrganizationTableData';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';
import type { Filter, Organization, SortBy } from '@spaces/graphql';
import { Table } from '@spaces/atoms/table';
import { useRecoilState, useRecoilValue } from 'recoil';
import { finderOrganizationsSearchTerms } from '../../../state';
import { mapGCliSearchTermsToFilterList } from '../../../utils/mapGCliSearchTerms';
import { finderOrganizationTableSortingState } from '../../../state/finderTables';
import { Building } from '@spaces/atoms/icons';

interface OrganizationListProps {
  preFilters?: Array<Filter>;
}

export const OrganizationList: React.FC<OrganizationListProps> = ({
  preFilters,
}: OrganizationListProps) => {
  const [page, setPagination] = useState(1);
  const [organizationsSearchTerms, setOrganizationsSearchTerms] =
    useRecoilState(finderOrganizationsSearchTerms);

  const { data, loading, fetchMore, variables, totalElements, refetchData } =
    useFinderOrganizationTableData(preFilters);
  const sortingState = useRecoilValue(finderOrganizationTableSortingState);

  const handleFilterResults = (searchTerms: any[]) => {
    setOrganizationsSearchTerms(searchTerms);
    setPagination(1);
    const sortBy: SortBy | undefined = sortingState.column
      ? {
          by: sortingState.column,
          direction: sortingState.direction,
          caseSensitive: false,
        }
      : undefined;
    let filters = mapGCliSearchTermsToFilterList(searchTerms, 'ORGANIZATION');
    if (preFilters) {
      filters = [...filters, ...preFilters];
    }
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 20,
        },
        where: {
          AND: filters,
        },
        sort: sortBy,
      },
    });
  };

  const [suggestions, setSuggestions] = useState<any[]>([]);
  const { data: gcliData, loading: gcliLoading, refetch } = useGCliSearch();

  useEffect(() => {
    if (!gcliLoading && gcliData) {
      setSuggestions(gcliData);
    }
  }, [gcliLoading, gcliData]);

  return (
    <>
      <div className={styles.inputSection}>
        <GCLIContextProvider
          label={'Organizations'}
          icon={<Building width={24} height={24} />}
          existingTerms={organizationsSearchTerms}
          loadSuggestions={(searchTerm: string) => {
            refetch && refetch({ limit: 5, keyword: searchTerm });
          }}
          loadingSuggestions={gcliLoading}
          suggestionsLoaded={suggestions}
          onItemsChange={handleFilterResults}
          selectedTermFormat={(item: any) => {
            if (item.type === 'STATE') {
              return item.data[0].value;
            }
            return item.display;
          }}
        >
          <GCLIInput />
        </GCLIContextProvider>
      </div>

      <Table<Organization>
        data={data}
        columns={organizationListColumns}
        isFetching={loading}
        totalItems={totalElements || 0}
        onFetchNextPage={() => {
          setPagination(page + 1);
          fetchMore({
            variables: {
              pagination: {
                limit: variables.pagination.limit,
                page: page + 1,
              },
            },
          });
        }}
      />
    </>
  );
};

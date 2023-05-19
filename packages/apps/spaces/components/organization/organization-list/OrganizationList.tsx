import React, { useEffect, useState } from 'react';
import styles from './organization-list.module.scss';
import { organizationListColumns } from './columns/OrganizationListColumns';
import { useFinderOrganizationTableData } from '@spaces/hooks/useFinderOrganizationTableData';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import SvgGlobe from '@spaces/atoms/icons/Globe';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';
import { Organization } from '../../../graphQL/__generated__/generated';
import { Table } from '@spaces/atoms/table';
import { useRecoilState } from 'recoil';
import {
  finderOrganizationsSearchTerms,
} from '../../../state';
import { mapGCliSearchTermsToFilterList } from '../../../utils/mapGCliSearchTerms';

export const OrganizationList: React.FC = () => {
  const [page, setPagination] = useState(1);

  const [organizationsSearchTerms, setOrganizationsSearchTerms] =
    useRecoilState(finderOrganizationsSearchTerms);

  const { data, loading, fetchMore, variables, totalElements } =
    useFinderOrganizationTableData(
      mapGCliSearchTermsToFilterList(organizationsSearchTerms, 'ORGANIZATION'),
    );

  const handleFilterResults = (searchTerms: any[]) => {
    setOrganizationsSearchTerms(searchTerms);
    setPagination(1);

    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 20,
        },
        where: {
          AND: mapGCliSearchTermsToFilterList(searchTerms, 'ORGANIZATION'),
        },
      },
    });
  };

  const [suggestions, setSuggestions] = useState<any[]>([]);
  const { data: gcliData, loading: gcliLoading, refetch } = useGCliSearch();

  useEffect(() => {
    if (!gcliLoading && gcliData) {
      setSuggestions(gcliData.map((item: any) => item.result));
    }
  }, [gcliLoading, gcliData]);

  return (
    <div style={{ padding: '1.2rem', height: '100%' }}>
      <div className={styles.inputSection}>
        <GCLIContextProvider
          label={'Organizations'}
          icon={<SvgGlobe />}
          inputPlaceholder={'Search by name or state'}
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
    </div>
  );
};

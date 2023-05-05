import React, {useEffect, useState} from 'react';
import styles from './organization-list.module.scss';
import {organizationListColumns} from './OrganizationListColumns';
import {useFinderOrganizationTableData} from '@spaces/hooks/useFinderOrganizationTableData';
import {useGCliSearch} from '@spaces/hooks/useGCliSearch';
import SvgGlobe from '@spaces/atoms/icons/Globe';
import {GCLIContextProvider, GCLIInput} from '@spaces/molecules/gCLI';
import {Filter, Organization} from "../../../graphQL/__generated__/generated";
import {Table} from "@spaces/atoms/table";

export const OrganizationList: React.FC = () => {
  const [page, setPagination] = useState(1);

  const {data, loading, fetchMore, variables, totalElements} = useFinderOrganizationTableData();

  const handleFilterResults = (searchTerms: any[]) => {
    setPagination(1);
    // setSearchTerm(value);
    const filters = [] as any[];
    searchTerms.forEach((item: any) => {
      if(item.type === 'STATE') {
        filters.push({
          filter: {
            property: 'REGION',
            operation: 'EQ',
            value: item.display,
          }
        });
        filters.push({
          filter: {
            property: 'REGION',
            operation: 'EQ',
            value: item.data[0].value,
          }
        });
      } else {
        filters.push({
          filter: {
            property: 'ORGANIZATION',
            operation: 'EQ',
            value: item.display,
          }
        });
      }
    });
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 60,
        },
        where: {
          AND: filters,
        } as Filter,
      },
    });
  };

  const [suggestions, setSuggestions] = useState<any[]>([]);
  const {
    data: gcliData,
    loading: gcliLoading,
    refetch,
  } = useGCliSearch();

  useEffect(() => {
    if (!gcliLoading && gcliData) {
      setSuggestions(gcliData.map((item: any) => item.result));
    }
  }, [gcliLoading, gcliData]);

  return (
      <div style={{padding: '1.2rem', height: '100%'}}>
        <div className={styles.inputSection}>
          <GCLIContextProvider
              label={'Organizations'}
              icon={<SvgGlobe/>}
              inputPlaceholder={
                'Search by name or state'
              }
              loadSuggestions={(searchTerm: string) => {
                refetch && refetch({limit: 5, keyword: searchTerm});
              }}
              loadingSuggestions={gcliLoading}
              suggestionsLoaded={suggestions}
              onItemsChange={handleFilterResults}
          >
            <GCLIInput/>
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

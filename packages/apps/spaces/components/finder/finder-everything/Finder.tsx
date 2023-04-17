import React, { useEffect, useState } from 'react';
import styles from './finder.module.scss';
import { useRecoilState } from 'recoil';
import { finderSearchTerm } from '../../../state';
import {
  DashboardViewItem,
  useFinderTableData,
} from '../../../hooks/useFinderTableData';
import { DebouncedInput, Table } from '../../ui-kit';
import { Search } from '../../ui-kit/atoms';
import { columns } from './Columns';
import { useRouter } from 'next/router';

export const Finder: React.FC = () => {
  const [page, setPagination] = useState(0);
  const [searchTerm, setSearchTerm] = useRecoilState(finderSearchTerm);
  const router = useRouter();

  const { data, loading, fetchMore, variables, totalElements } =
    useFinderTableData({
      pagination: {
        page: 1,
        limit: 60,
      },
      searchTerm,
    });

  useEffect(() => {
    // Parse the search query from the URL
    const { query } = router;
    const searchQueryFromUrl = (query.q as string) || '';

    // Update the search input with the query from the URL
    setSearchTerm(searchQueryFromUrl);

    // Fetch the list based on the search query

    setPagination(1);
    setSearchTerm(searchQueryFromUrl);
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 60,
        },
        searchTerm: searchQueryFromUrl,
      },
    });
  }, [router]);

  useEffect(() => {
    // revisit later
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 60,
        },
        searchTerm: '',
      },
    });
  }, []);

  useEffect(() => {
    const param = new URL(window.location.href).searchParams.get('q');
    if (param) console.log(param);
  }, [searchTerm]);

  const handleFilterResults = (value: string) => {
    setPagination(1);
    setSearchTerm(value);
    router.push(`/?q=${value}`);
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 60,
        },
        searchTerm: value,
      },
    });
  };

  return (
    <div style={{ padding: '1.2rem', height: '100%' }}>
      <div className={styles.inputSection}>
        <h1> Everything </h1>

        <DebouncedInput
          // todo temporary
          minLength={1}
          onChange={(event: any) => handleFilterResults(event.target.value)}
          placeholder={'Search organizations, contacts, locations...'}
          value={searchTerm}
        >
          <Search />
        </DebouncedInput>
      </div>

      <Table<DashboardViewItem>
        data={data}
        columns={columns}
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
              searchTerm,
            },
          });
        }}
      />
    </div>
  );
};

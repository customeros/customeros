import React, { useEffect, useState } from 'react';
import styles from './finder-contact.module.scss';
import { useRecoilState } from 'recoil';
import { finderSearchTerm } from '../../../state';
import { Table } from '@spaces/atoms/table';
import { DebouncedInput } from '@spaces/atoms/input';
import Search from '@spaces/atoms/icons/Search';
import { finderContactColumns } from './FinderContactColumns';
import { useRouter } from 'next/router';
import { Contact } from '../../../graphQL/__generated__/generated';
import { useFinderContactTableData } from '@spaces/hooks/useFinderContactTableData';

export const FinderContact: React.FC = () => {
  const [page, setPagination] = useState(0);
  const [searchTerm, setSearchTerm] = useRecoilState(finderSearchTerm);
  const router = useRouter();

  const { data, loading, fetchMore, variables, totalElements } =
    useFinderContactTableData({
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
        // searchTerm: searchQueryFromUrl, TODO fix this
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
        // searchTerm: '', TODO fix this
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
        // searchTerm: value,  TODO fix this
      },
    });
  };

  return (
    <div style={{ padding: '1.2rem', height: '100%' }}>
      <div className={styles.inputSection}>
        <h1> Contacts </h1>

        <DebouncedInput
          // todo temporary
          minLength={1}
          onChange={(event: any) => handleFilterResults(event.target.value)}
          placeholder={'Search contacts, locations...'}
          value={searchTerm}
        >
          <Search />
        </DebouncedInput>
      </div>

      <Table<Contact>
        data={data}
        columns={finderContactColumns}
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
              // searchTerm,  TODO fix this
            },
          });
        }}
      />
    </div>
  );
};

import React, { useState } from 'react';
import { Table, DebouncedInput } from '../ui-kit';
import { columns } from './finder-table/Columns';
import styles from './finder.module.scss';
import { DashboardViewItem } from '../../graphQL/__generated__/generated';
import { Search } from '../ui-kit/atoms';
import { useFinderTableData } from '../../hooks/useFinderTableData';
import { useRecoilState } from 'recoil';
import { finderSearchTerm } from '../../state';

export const Finder: React.FC = () => {
  const [page, setPagination] = useState(0);
  const [searchTerm, setSearchTerm] = useRecoilState(finderSearchTerm);

  const { data, loading, fetchMore, variables, totalElements } =
    useFinderTableData({
      pagination: {
        page: 0,
        limit: 60,
      },
      searchTerm,
    });

  const handleFilterResults = (value: string) => {
    setPagination(0);
    setSearchTerm(value);
    fetchMore({
      variables: {
        pagination: {
          page: 0,
          limit: 10,
        },
        searchTerm: value,
      },
    });
  };
  return (
    <div style={{ padding: '36px', height: '100%' }}>
      <div className={styles.inputSection}>
        <h1> Everything </h1>

        <DebouncedInput
          // todo temporary
          minLength={1}
          onChange={(event) => handleFilterResults(event.target.value)}
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

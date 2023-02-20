import React, { useState } from 'react';
import { useGetDashboardData } from './graphQL/useGetDashboardData';
import { Table, DebouncedInput } from '../ui-kit';
import { columns } from './Columns';
import styles from './finder.module.scss';
import { DashboardViewItem } from '../../graphQL/types';

export const Finder: React.FC = () => {
  const [page, setPagination] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  const { data, loading, fetchMore, variables, totalElements } =
    useGetDashboardData({
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
        <h1> All Contacts </h1>

        <DebouncedInput
          onChange={(event) => handleFilterResults(event.target.value)}
        />
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

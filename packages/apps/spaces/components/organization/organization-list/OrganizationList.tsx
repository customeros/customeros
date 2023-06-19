import React, { useEffect, useState } from 'react';
import styles from './organization-list.module.scss';
import { columns } from './OrganizationListColumns';
import { useFinderOrganizationTableData } from '@spaces/hooks/useFinderOrganizationTableData';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';
import { Button } from '@spaces/atoms/button';
import {
  useCreateOrganization,
  useMergeOrganizations,
} from '@spaces/hooks/useOrganization';
import {
  SortingDirection,
  type Filter,
  type Organization,
  type SortBy,
} from '@spaces/graphql';
import {
  Table,
  RowSelectionState,
  SortingState,
  TableInstance,
} from '@spaces/ui/presentation/Table/Table';
import { TActions } from '@spaces/ui/presentation/Table/TActions';
import { useRecoilState } from 'recoil';
import { finderOrganizationsSearchTerms } from '../../../state';
import { mapGCliSearchTermsToFilterList } from '../../../utils/mapGCliSearchTerms';
import { useRouter } from 'next/router';
import { IconButton } from '@spaces/atoms/icon-button';
import { Check } from '@spaces/atoms/icons';

const skeletonData = Array.from({ length: 5 }).map(
  () => ({}),
) as Organization[];

interface OrganizationListProps {
  preFilters?: Array<Filter>;
  label: string;
  icon: React.ReactNode;
}

export const OrganizationList: React.FC<OrganizationListProps> = ({
  preFilters,
  label,
  icon,
}: OrganizationListProps) => {
  const [page, setPagination] = useState(1);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [enableSelection, setEnableSelection] = useState(false);
  const [selection, setSelection] = useState<RowSelectionState>({});

  const sortBy: SortBy | undefined = (() => {
    if (!sorting.length) return;
    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  })();

  const { push } = useRouter();
  const { onMergeOrganizations } = useMergeOrganizations();
  const { onCreateOrganization } = useCreateOrganization();

  const [organizationsSearchTerms, setOrganizationsSearchTerms] =
    useRecoilState(finderOrganizationsSearchTerms);
  const { data, loading, fetchMore, variables, totalElements } =
    useFinderOrganizationTableData(preFilters, sortBy);

  const handleFilterResults = (searchTerms: any[]) => {
    setOrganizationsSearchTerms(searchTerms);
    setPagination(1);

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

  const tableActions = [
    {
      label: 'Add organization',
      command: async () => {
        const newOrganization = await onCreateOrganization({ name: '' });
        if (newOrganization?.id) {
          push(`/organization/${newOrganization?.id}`);
        }
      },
    },
    {
      label: 'Merge organizations',
      command: () => {
        setEnableSelection(true);
      },
    },
  ];

  const handleFetchMore = () => {
    setPagination(page + 1);
    fetchMore({
      variables: {
        pagination: {
          limit: variables.pagination.limit,
          page: page + 1,
        },
      },
    });
  };

  const handleMergeOrganizations = (table: TableInstance<Organization>) => {
    const organizationIds = Object.keys(selection)
      .map((key) => data?.[Number(key)]?.id)
      .filter(Boolean) as string[];

    const primaryId = organizationIds[0];
    const mergeIds = organizationIds.slice(1);

    onMergeOrganizations({
      primaryOrganizationId: primaryId,
      mergedOrganizationIds: mergeIds,
    });
    table.resetRowSelection();
  };

  useEffect(() => {
    if (!gcliLoading && gcliData) {
      setSuggestions(gcliData);
    }
  }, [gcliLoading, gcliData]);

  return (
    <>
      <div className={styles.inputSection}>
        <GCLIContextProvider
          label={label}
          icon={icon}
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
        data={loading ? skeletonData : data ?? []}
        columns={columns}
        sorting={sorting}
        enableTableActions
        isLoading={loading}
        selection={selection}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        totalItems={totalElements || 0}
        onSelectionChange={setSelection}
        enableRowSelection={enableSelection}
        renderTableActions={(ref, table) => {
          if (enableSelection) {
            if (Object.keys(selection).length > 1) {
              return (
                <div ref={ref}>
                  <Button
                    mode='primary'
                    onClick={() => handleMergeOrganizations(table)}
                  >
                    Merge
                  </Button>
                </div>
              );
            }
            return (
              <div
                ref={ref}
                style={{
                  display: 'flex',
                  height: '100%',
                  alignItems: 'center',
                }}
              >
                <IconButton
                  mode='secondary'
                  style={{ padding: 0 }}
                  onClick={() => {
                    setEnableSelection(false);
                    table.resetRowSelection();
                  }}
                  label='Done'
                  icon={<Check height={24} width={24} />}
                />
              </div>
            );
          }
          return <TActions ref={ref} model={tableActions} />;
        }}
      />
    </>
  );
};

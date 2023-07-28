'use client';

import React, { useEffect, useMemo, useState, lazy, Suspense } from 'react';
import styles from './organization-list.module.scss';
import { columns } from './OrganizationListColumns';
import { useFinderOrganizationTableData } from '@spaces/hooks/useFinderOrganizationTableData';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';

import {
  Table,
  RowSelectionState,
  TableInstance,
  SortingState,
} from '@ui/presentation/Table';

import {
  useArchiveOrganizations,
  useCreateOrganization,
  useMergeOrganizations,
} from '@spaces/hooks/useOrganization';
import {
  SortingDirection,
  type Filter,
  type Organization,
  type SortBy,
} from '@spaces/graphql';

import { useRecoilState } from 'recoil';
import { finderOrganizationsSearchTerms } from '../../../state';
import { mapGCliSearchTermsToFilterList } from '@spaces/utils/mapGCliSearchTerms';
import { useRouter } from 'next/router';
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@ui/presentation/Modal';
import { Button } from '@ui/form/Button';
import { useDisclosure } from '@chakra-ui/react-use-disclosure';

const OrganizationListActions = lazy(() => import('./OrganizationListActions'));

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
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [tableInstance, setTableInstance] =
    useState<TableInstance<Organization> | null>(null);
  const [page, setPagination] = useState(1);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [enableSelection, setEnableSelection] = useState(false);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [suggestions, setSuggestions] = useState<any[]>([]);
  const { data: gcliData, loading: gcliLoading, refetch } = useGCliSearch();
  const sortBy: SortBy | undefined = useMemo(() => {
    setPagination(1);
    if (!sorting.length) return;
    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  }, [sorting]);

  const { push } = useRouter();

  const { onMergeOrganizations } = useMergeOrganizations();
  const { onArchiveOrganization } = useArchiveOrganizations();
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

  const handleCreateOrganization = async () => {
    const newOrganization = await onCreateOrganization({ name: '' });
    if (newOrganization?.id) {
      push(`/organizations/${newOrganization?.id}`);
    }
  };

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
  const handleArchiveOrganizations = () => {
    const organizationIds = Object.keys(selection)
      .map((key) => data?.[Number(key)]?.id)
      .filter(Boolean) as string[];
    onArchiveOrganization({
      ids: organizationIds,
    });
    onClose();
    tableInstance?.resetRowSelection();
    setEnableSelection(false);
    setTableInstance(null);
  };
  const handleCancelRemoveOrganizations = () => {
    onClose();
    tableInstance?.resetRowSelection();
    setEnableSelection(false);
    setTableInstance(null);
  };
  const handleOpenConfirmationModal = (table: TableInstance<Organization>) => {
    setTableInstance(table);
    onOpen();
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
      <Modal onClose={onClose} isOpen={isOpen} size='xl'>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>
            Are you sure you want to delete those organizations?
            <ModalCloseButton onClick={handleCancelRemoveOrganizations} mt={2} />

          </ModalHeader>
          <ModalFooter>
            <Button onClick={handleCancelRemoveOrganizations} mr={2}>
              Cancel
            </Button>

            <Button onClick={handleArchiveOrganizations} colorScheme='red'>
              Delete
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Table<Organization>
        data={data ?? []}
        columns={columns}
        sorting={sorting}
        enableTableActions
        isLoading={loading}
        selection={selection}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        totalItems={totalElements ?? 0}
        onSelectionChange={setSelection}
        enableRowSelection={enableSelection}
        renderTableActions={(table) => (
          <Suspense fallback={<div />}>
            <OrganizationListActions
              table={table}
              selection={selection}
              isSelectionEnabled={enableSelection}
              toggleSelection={setEnableSelection}
              onCreateOrganization={handleCreateOrganization}
              onMergeOrganizations={handleMergeOrganizations}
              onArchiveOrganizations={handleOpenConfirmationModal}
            />
          </Suspense>
        )}
      />
    </>
  );
};

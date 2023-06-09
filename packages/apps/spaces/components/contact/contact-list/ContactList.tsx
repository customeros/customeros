import React, { useEffect, useState } from 'react';
import styles from './contact-list.module.scss';
import { contactListColumns } from './columns/ContactListColumns';
import { Contact } from '../../../graphQL/__generated__/generated';
import { useFinderContactTableData } from '@spaces/hooks/useFinderContactTableData';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';
import { Table } from '@spaces/atoms/table';
import { useRecoilState } from 'recoil';
import { finderContactsGridData } from '../../../state';
import { mapGCliSearchTermsToFilterList } from '../../../utils/mapGCliSearchTerms';
import { User } from '@spaces/atoms/icons';

export const ContactList: React.FC = () => {
  const [page, setPagination] = useState(1);

  const [contactsGridData, setContactsGridData] = useRecoilState(
    finderContactsGridData,
  );

  const [suggestions, setSuggestions] = useState<any[]>([]);
  const { data: gcliData, loading: gcliLoading, refetch } = useGCliSearch();
  const mapSortBy = contactsGridData.sortBy.column
    ? {
        by: contactsGridData.sortBy.column,
        direction: contactsGridData.sortBy.direction,
        caseSensitive: false,
      }
    : undefined;
  const { data, loading, fetchMore, variables, totalElements } =
    useFinderContactTableData(
      mapGCliSearchTermsToFilterList(contactsGridData.searchTerms, 'CONTACT'),
      mapSortBy,
    );

  const handleFilterResults = (searchTerms: any[]) => {
    setContactsGridData((prevState: any) => ({
      ...prevState,
      searchTerms,
    }));
    setPagination(1);
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 20,
        },
        where: { AND: mapGCliSearchTermsToFilterList(searchTerms, 'CONTACT') },
        sort: mapSortBy,
      },
    });
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
          label={'Contacts'}
          icon={<User width={24} height={24} />}
          existingTerms={contactsGridData.searchTerms}
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

      <Table<Contact>
        data={data}
        columns={contactListColumns}
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

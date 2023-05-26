import React, { useEffect, useState } from 'react';
import styles from './contact-list.module.scss';
import { contactListColumns } from './columns/ContactListColumns';
import { Contact, SortBy } from '../../../graphQL/__generated__/generated';
import { useFinderContactTableData } from '@spaces/hooks/useFinderContactTableData';
import SvgGlobe from '@spaces/atoms/icons/Globe';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';
import { Table } from '@spaces/atoms/table';
import { useRecoilState, useRecoilValue } from 'recoil';
import { finderContactsSearchTerms } from '../../../state';
import { mapGCliSearchTermsToFilterList } from '../../../utils/mapGCliSearchTerms';
import { finderContactTableSortingState } from '../../../state/finderTables';
import { User } from '@spaces/atoms/icons';

export const ContactList: React.FC = () => {
  const [page, setPagination] = useState(1);

  const [contactsSearchTerms, setContactsSearchTerms] = useRecoilState(
    finderContactsSearchTerms,
  );
  const [suggestions, setSuggestions] = useState<any[]>([]);
  const { data: gcliData, loading: gcliLoading, refetch } = useGCliSearch();
  const { data, loading, fetchMore, variables, totalElements } =
    useFinderContactTableData(
      mapGCliSearchTermsToFilterList(contactsSearchTerms, 'CONTACT'),
    );
  const sortingState = useRecoilValue(finderContactTableSortingState);

  const handleFilterResults = (searchTerms: any[]) => {
    setContactsSearchTerms(searchTerms);
    setPagination(1);
    const sortBy: SortBy | undefined = sortingState.column
      ? {
          by: sortingState.column,
          direction: sortingState.direction,
          caseSensitive: false,
        }
      : undefined;
    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 20,
        },
        where: { AND: mapGCliSearchTermsToFilterList(searchTerms, 'CONTACT') },
        sort: sortBy,
      },
    });
  };

  useEffect(() => {
    if (!gcliLoading && gcliData) {
      setSuggestions(gcliData.map((item: any) => item.result));
    }
  }, [gcliLoading, gcliData]);

  return (
    <>
      <div className={styles.inputSection}>
        <GCLIContextProvider
          label={'Contacts'}
          icon={<User width={24} height={24} />}
          inputPlaceholder={'search and filter'}
          existingTerms={contactsSearchTerms}
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

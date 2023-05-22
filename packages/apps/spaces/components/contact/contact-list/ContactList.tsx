import React, { useEffect, useState } from 'react';
import styles from './contact-list.module.scss';
import { contactListColumns } from './columns/ContactListColumns';
import { Contact } from '../../../graphQL/__generated__/generated';
import { useFinderContactTableData } from '@spaces/hooks/useFinderContactTableData';
import SvgGlobe from '@spaces/atoms/icons/Globe';
import { useGCliSearch } from '@spaces/hooks/useGCliSearch';
import { GCLIContextProvider, GCLIInput } from '@spaces/molecules/gCLI';
import { Table } from '@spaces/atoms/table';
import { useRecoilState } from 'recoil';
import { finderContactsSearchTerms } from '../../../state';
import { mapGCliSearchTermsToFilterList } from '../../../utils/mapGCliSearchTerms';

export const ContactList: React.FC = () => {
  const [page, setPagination] = useState(1);

  const [contactsSearchTerms, setContactsSearchTerms] = useRecoilState(
    finderContactsSearchTerms,
  );
  const { data, loading, fetchMore, variables, totalElements } =
    useFinderContactTableData(
      mapGCliSearchTermsToFilterList(contactsSearchTerms, 'CONTACT'),
    );

  const handleFilterResults = (searchTerms: any[]) => {
    setContactsSearchTerms(searchTerms);
    setPagination(1);

    fetchMore({
      variables: {
        pagination: {
          page: 1,
          limit: 20,
        },
        where: { AND: mapGCliSearchTermsToFilterList(searchTerms, 'CONTACT') },
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
    <>
      <div className={styles.inputSection}>
        <GCLIContextProvider
          label={'Contacts'}
          icon={<SvgGlobe />}
          inputPlaceholder={'Search by name or state'}
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

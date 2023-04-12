import React, { useEffect, useMemo, useState } from 'react';
import {
  ComparisonOperator,
  useContactMentionSuggestionsList,
} from '../../../../../hooks/useContactList';
import {
  Contact,
  Organization,
} from '../../../../../graphQL/__generated__/generated';
import {
  MentionAtomState,
  MentionAtomNodeAttributes,
  MentionAtomPopupComponent,
} from '@remirror/react';
import { useOrganizationMentionSuggestionsList } from '../../../../../hooks/useOrganizations';

interface MentionProps<
  UserData extends MentionAtomNodeAttributes = MentionAtomNodeAttributes,
> {
  users?: UserData[];
  tags?: string[];
}

export const Mention = () => {
  const [mentionState, setMentionState] = useState<MentionAtomState | null>();
  const { onLoadContactMentionSuggestionsList } =
    useContactMentionSuggestionsList();
  const { onLoadOrganizationMentionSuggestionsList } =
    useOrganizationMentionSuggestionsList();
  const [filteredContacts, setFilteredContacts] = useState<
    Array<{ id: string; label: string }>
  >([]);
  const [filteredOrganizations, setFilteredOrganizations] = useState<
    Array<{ id: string; label: string }>
  >([]);

  const getContactSuggestions = async (filter: string) => {
    const response = await onLoadContactMentionSuggestionsList({
      variables: {
        pagination: { page: 0, limit: 10 },
        where: {
          OR: [
            {
              filter: {
                property: 'FIRST_NAME',
                value: filter,
                operation: ComparisonOperator.Contains,
              },
            },
            {
              filter: {
                property: 'LAST_NAME',
                value: filter,
                operation: ComparisonOperator.Contains,
              },
            },
          ],
        },
      },
    });
    if (response?.data) {
      const options = response?.data?.contacts?.content.map((e: Contact) => ({
        label: `${e.firstName} ${e.lastName}`,
        id: e.id,
      }));
      setFilteredContacts(options || []);
    }
  };

  const getOrganizationSuggestions = async (filter: string) => {
    const response = await onLoadOrganizationMentionSuggestionsList({
      variables: {
        pagination: { page: 0, limit: 10 },
        where: {
          filter: {
            property: 'NAME',
            value: filter,
            operation: ComparisonOperator.Contains,
          },
        },
      },
    });
    if (response?.data) {
      const options = response?.data?.organizations?.content.map(
        (e: Organization) => ({
          label: `${e.name}`,
          id: e.id,
        }),
      );
      setFilteredOrganizations(options || []);
    }
  };

  useEffect(() => {
    if (mentionState?.name === 'at' && mentionState?.query?.full) {
      getContactSuggestions(mentionState.query.full);
    }

    if (mentionState?.name !== 'at' && mentionState?.query?.full) {
      getOrganizationSuggestions(mentionState.query.full);
    }
  }, [mentionState?.query?.full, mentionState?.name]);

  const items = useMemo(() => {
    if (!mentionState) {
      return [];
    }

    const allItems =
      mentionState.name === 'at' ? filteredContacts : filteredOrganizations;

    if (!allItems) {
      return [];
    }

    const query = mentionState.query.full.toLowerCase() ?? '';
    return allItems
      .filter((item) => item?.label.toLowerCase().includes(query))
      .sort();
  }, [mentionState, filteredOrganizations, filteredContacts]);

  return <MentionAtomPopupComponent onChange={setMentionState} items={items} />;
};

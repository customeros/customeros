import React, { useEffect, useState } from 'react';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { Autocomplete } from '@spaces/atoms/new-autocomplete';
import { SuggestionItem } from '@spaces/atoms/new-autocomplete/Autocomplete';
import { useRecoilValue } from 'recoil';
import { ownerListData } from '../../../../state/userData';
import Fuse from 'fuse.js';

interface OrganizationOwnerProps {
  id: string;
  editMode: boolean;
  switchEditMode?: () => void;
  owner: any;
}

export const OrganizationOwnerAutocomplete: React.FC<
  OrganizationOwnerProps
> = ({ id, editMode, switchEditMode, owner }) => {
  const [inputValue, setInputValue] = React.useState<string>('');

  const { onLinkOrganizationOwner, saving } = useLinkOrganizationOwner({
    organizationId: id,
  });
  const { onUnlinkOrganizationOwner } = useUnlinkOrganizationOwner({
    organizationId: id,
  });

  useEffect(() => {
    setInputValue(owner ? owner.firstName + ' ' + owner.lastName : '');
  }, [owner]);

  const [ownerListMatch, setOwnerListMatch] = useState<Array<SuggestionItem>>(
    [],
  );
  const [ownerListFuzzy, setOwnerListFuzzy] = useState<Array<SuggestionItem>>(
    [],
  );

  const [loadingOwnerSuggestions, setLoadingOwnerSuggestions] =
    useState<boolean>(false);
  const { ownerList } = useRecoilValue(ownerListData);
  const fuse = new Fuse(ownerList, { keys: ['firstName', 'lastName'] });

  const mapOwnerToSuggestionItem = (owner: any) => {
    return { label: `${owner.firstName} ${owner.lastName}`, value: owner.id };
  };

  const handleChangeOwner = ({
    value,
    label,
  }: {
    value: string;
    label: string;
  }) => {
    onLinkOrganizationOwner({ userId: value, name: label }).then(() => {
      switchEditMode && switchEditMode();
    });
  };

  if (!editMode) {
    return (
      <div
        tabIndex={0}
        role='button'
        style={{
          cursor: 'pointer',
          overflowX: 'hidden',
          textOverflow: 'ellipsis',
          border: '1px solid transparent',
        }}
        onDoubleClick={switchEditMode}
        onKeyDown={(e) => {
          e.key === 'Enter' && switchEditMode && switchEditMode();
        }}
      >
        {inputValue || <span style={{ color: '#ccc' }}>Owner </span>}
      </div>
    );
  }

  const filterUsers = (users: any[], filter: string): string[] => {
    const filterWords = filter.toLowerCase().split(' ');

    return users.filter((owner) => {
      const firstName = owner.firstName.toLowerCase();
      const lastName = owner.lastName.toLowerCase();

      if (filterWords.length === 1) {
        // Match users with either name or surname containing the filter string
        return (
          firstName.includes(filter.toLowerCase()) ||
          lastName.includes(filter.toLowerCase())
        );
      } else if (filterWords.length > 1) {
        // Match users with name or surname equals the first filter word
        // and name or surname starts with the second filter word
        const filterStartsWithFirstName =
          firstName.includes(filterWords[0]) &&
          lastName.startsWith(filterWords[1]);
        const filterStartsWithLastName =
          lastName.includes(filterWords[0]) &&
          firstName.startsWith(filterWords[1]);

        return filterStartsWithFirstName || filterStartsWithLastName;
      }

      return false;
    });
  };

  const searchOwners = (filter: string) => {
    setLoadingOwnerSuggestions(true);
    const ownersContains = filterUsers(ownerList, filter);

    if (ownersContains.length > 0) {
      setOwnerListFuzzy([]);
      setOwnerListMatch(ownersContains.map(mapOwnerToSuggestionItem));
    } else {
      setOwnerListMatch([]);
      const ownersFuzzySearch = fuse.search(filter);
      if (ownersFuzzySearch.length > 0) {
        setOwnerListFuzzy(
          ownersFuzzySearch.map((p) => p.item).map(mapOwnerToSuggestionItem),
        );
      } else {
        setOwnerListFuzzy(ownerList.map(mapOwnerToSuggestionItem));
      }
    }
    setLoadingOwnerSuggestions(false);
  };

  return (
    <Autocomplete
      loading={false}
      mode='full-width'
      editable
      initialValue={inputValue}
      suggestionsMatch={ownerListMatch}
      suggestionsFuzzyMatch={ownerListFuzzy}
      onDoubleClick={switchEditMode}
      onChange={(e: any) => {
        if (e?.value) {
          handleChangeOwner(e);
        }
      }}
      saving={saving}
      onSearch={searchOwners}
      onClearInput={() => {
        if (owner) {
          onUnlinkOrganizationOwner();
        }
      }}
      placeholder='Owner'
    />
  );
};

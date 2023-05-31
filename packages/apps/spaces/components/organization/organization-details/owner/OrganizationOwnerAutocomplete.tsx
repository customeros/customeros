import React, { useEffect, useState } from 'react';
import { useOrganizationOwner } from '@spaces/hooks/useOrganizationOwner';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUserSuggestionsList } from '@spaces/hooks/useUser';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { Autocomplete } from '@spaces/atoms/new-autocomplete';
import {SuggestionItem} from "@spaces/atoms/new-autocomplete/Autocomplete";

interface OrganizationOwnerProps {
  id: string;
  editMode: boolean;
  switchEditMode?: () => void;
}

export const OrganizationOwnerAutocomplete: React.FC<
  OrganizationOwnerProps
> = ({ id, editMode, switchEditMode }) => {
  const [userSuggestions, setUserSuggestions] = useState<Array<SuggestionItem>>([]);
  const [userId, setUserId] = React.useState<string>('');
  const [inputValue, setInputValue] = React.useState<string>('');

  const { data, loading, error } = useOrganizationOwner({ id });
  const { getUsersSuggestions } = useUserSuggestionsList();
  const { onLinkOrganizationOwner } = useLinkOrganizationOwner({
    organizationId: id,
    userId,
  });
  const { onUnlinkOrganizationOwner } = useUnlinkOrganizationOwner({
    organizationId: id,
  });

  useEffect(() => {
    if (!loading && data) {
      setInputValue(
        data.owner ? data.owner.firstName + ' ' + data.owner.lastName : '',
      );
    }
  }, [data, loading]);

  useEffect(() => {
    if (userId) {
      onLinkOrganizationOwner().then(() => {
        switchEditMode && switchEditMode();
        setUserId('')
      });
    }
  }, [userId]);

  return (
    <>
      <Autocomplete
        mode='full-width'
        editable={editMode}
        initialValue={inputValue}
        suggestions={userSuggestions}
        onDoubleClick={() => {
          !editMode && switchEditMode && switchEditMode();
        }}
        onChange={(e: any) => {
          setUserId(e.value);
        }}
        loading={loading}
        onSearch={(filter: string) =>
          getUsersSuggestions(filter).then((options) => setUserSuggestions(options) )
        }
        onClearInput={() => {
          if (data?.owner) {
            onUnlinkOrganizationOwner();
          }
        }}
        placeholder='Search for a user'
      />
    </>
  );
};

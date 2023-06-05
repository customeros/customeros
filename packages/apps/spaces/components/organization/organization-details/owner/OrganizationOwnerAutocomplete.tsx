import React, { useEffect, useState } from 'react';
import { useOrganizationOwner } from '@spaces/hooks/useOrganizationOwner';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUserSuggestionsList } from '@spaces/hooks/useUser';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { Autocomplete } from '@spaces/atoms/new-autocomplete';
import { SuggestionItem } from '@spaces/atoms/new-autocomplete/Autocomplete';

interface OrganizationOwnerProps {
  id: string;
  editMode: boolean;
  switchEditMode?: () => void;
}

export const OrganizationOwnerAutocomplete: React.FC<
  OrganizationOwnerProps
> = ({ id, editMode, switchEditMode }) => {
  const [userSuggestions, setUserSuggestions] = useState<Array<SuggestionItem>>(
    [],
  );
  const [inputValue, setInputValue] = React.useState<string>('');

  const { data, loading, error } = useOrganizationOwner({ id });
  const { getUsersSuggestions } = useUserSuggestionsList();
  const { onLinkOrganizationOwner, saving } = useLinkOrganizationOwner({
    organizationId: id,
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
        style={{ marginTop: 8, cursor: 'pointer' }}
        onDoubleClick={switchEditMode}
        onKeyDown={(e) => {
          e.key === 'Enter' && switchEditMode && switchEditMode();
        }}
      >
        {inputValue || <span>Owner </span>}
      </div>
    );
  }

  return (
    <Autocomplete
      mode='full-width'
      editable
      initialValue={inputValue}
      suggestions={userSuggestions}
      onDoubleClick={switchEditMode}
      onChange={(e: any) => {
        if (e?.value) {
          handleChangeOwner(e);
        }
      }}
      loading={loading}
      saving={saving}
      onSearch={(filter: string) =>
        getUsersSuggestions(filter).then((options) =>
          setUserSuggestions(options),
        )
      }
      onClearInput={() => {
        if (data?.owner) {
          onUnlinkOrganizationOwner();
        }
      }}
      placeholder='Owner'
    />
  );
};

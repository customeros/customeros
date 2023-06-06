import React, { useEffect, useState } from 'react';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUserSuggestionsList } from '@spaces/hooks/useUser';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { Autocomplete } from '@spaces/atoms/new-autocomplete';
import { SuggestionItem } from '@spaces/atoms/new-autocomplete/Autocomplete';

interface OrganizationOwnerProps {
  id: string;
  editMode: boolean;
  switchEditMode?: () => void;
  owner: any;
}

export const OrganizationOwnerAutocomplete: React.FC<
  OrganizationOwnerProps
> = ({ id, editMode, switchEditMode, owner }) => {
  const [userSuggestions, setUserSuggestions] = useState<Array<SuggestionItem>>(
    [],
  );
  const [inputValue, setInputValue] = React.useState<string>('');

  const { getUsersSuggestions } = useUserSuggestionsList();
  const { onLinkOrganizationOwner, saving } = useLinkOrganizationOwner({
    organizationId: id,
  });
  const { onUnlinkOrganizationOwner } = useUnlinkOrganizationOwner({
    organizationId: id,
  });

  useEffect(() => {
    setInputValue(owner ? owner.firstName + ' ' + owner.lastName : '');
  }, [owner]);

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

  return (
    <Autocomplete
      loading={false}
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
      saving={saving}
      onSearch={(filter: string) =>
        getUsersSuggestions(filter).then((options) => {
          setUserSuggestions(options);
        })
      }
      onClearInput={() => {
        if (owner) {
          onUnlinkOrganizationOwner();
        }
      }}
      placeholder='Owner'
    />
  );
};

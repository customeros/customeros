import React, { useEffect, useState } from 'react';
import { useOrganizationOwner } from '@spaces/hooks/useOrganizationOwner';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import styles from './organization-owner.module.scss';
import { DebouncedAutocomplete } from '@spaces/atoms/autocomplete';
import { useUserSuggestionsList } from '@spaces/hooks/useUser';
import { SearchPlus } from '@spaces/atoms/icons';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';

interface OrganizationOwnerProps {
  id: string;
}

export const OrganizationOwner: React.FC<OrganizationOwnerProps> = ({ id }) => {
  const [{ isEditMode }] = useRecoilState(organizationDetailsEdit);

  const [userSuggestions, setUserSuggestions] = useState<Array<any>>([]);
  const [userId, setUserId] = React.useState<string>('');
  const [inputValue, setInputValue] = React.useState<string | undefined>(
    undefined,
  );

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
      onLinkOrganizationOwner();
    }
  }, [userId]);

  if (loading) return null;
  if (error) {
    return (
      <div>Sorry looks like there was an error during loading the owner</div>
    );
  }

  if (!isEditMode) {
    return (
      <article className={styles.owner_section}>
        <h1 className={styles.owner_header}>Owner</h1>
        {!data?.owner && (
          <div className={styles.owner}>This company has no owner</div>
        )}
        {data?.owner && (
          <div className={styles.owner}>
            {data.owner ? data.owner.firstName + ' ' + data.owner.lastName : ''}
          </div>
        )}
      </article>
    );
  } else {
    return (
      <article className={styles.owner_section}>
        <h1 className={styles.owner_header}>Owner</h1>
        <div className={styles.owner_input}>
          <SearchPlus height={14} />

          {(inputValue || inputValue === '') && (
            <DebouncedAutocomplete
              mode='fit-content'
              editable={true}
              initialValue={inputValue}
              suggestions={userSuggestions}
              onChange={(e) => {
                setUserId(e.value);
              }}
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
              placeholder='Search for a user'
              newItemLabel=''
            />
          )}
        </div>
      </article>
    );
  }
};

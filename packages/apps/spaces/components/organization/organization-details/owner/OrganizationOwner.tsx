import React from 'react';
import { useOrganizationOwner } from '@spaces/hooks/useOrganizationOwner';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import styles from './organization-owner.module.scss';
import { OrganizationOwnerAutocomplete } from '@spaces/organization/organization-details/owner/OrganizationOwnerAutocomplete';

interface OrganizationOwnerProps {
  id: string;
}

export const OrganizationOwner: React.FC<OrganizationOwnerProps> = ({ id }) => {
  const [{ isEditMode }] = useRecoilState(organizationDetailsEdit);

  const { data, loading, error } = useOrganizationOwner({ id });

  return (
    <article className={styles.owner_section}>
      <h1 className={styles.owner_header}>Owner</h1>
      {!data?.owner && !isEditMode && (
        <div className={styles.owner}>This company has no owner</div>
      )}
      {data?.owner && !isEditMode && (
        <div className={styles.owner}>
          {data.owner ? data.owner.firstName + ' ' + data.owner.lastName : ''}
        </div>
      )}

      {isEditMode && (
        <OrganizationOwnerAutocomplete id={id} editMode={isEditMode} />
      )}
    </article>
  );
};

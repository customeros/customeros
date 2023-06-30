import React from 'react';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import styles from './organization-owner.module.scss';
import { OrganizationOwnerAutocomplete } from '@spaces/organization/organization-details/owner/OrganizationOwnerAutocomplete';
import { User } from '@spaces/graphql';
type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;

interface OrganizationOwnerProps {
  id: string;
  owner?: Owner;
}

export const OrganizationOwner: React.FC<OrganizationOwnerProps> = ({
  id,
  owner,
}) => {
  const [{ isEditMode }] = useRecoilState(organizationDetailsEdit);

  return (
    <article className={styles.owner_section}>
      <h1 className={styles.owner_header}>Owner</h1>
      {!owner && !isEditMode && (
        <div className={styles.owner}>This company has no owner</div>
      )}
      {owner && !isEditMode && (
        <div className={styles.owner}>
          {owner ? owner.firstName + ' ' + owner.lastName : ''}
        </div>
      )}

      {isEditMode && (
        <div style={{ position: 'relative' }}>
          <OrganizationOwnerAutocomplete id={id} owner={owner} />
        </div>
      )}
    </article>
  );
};

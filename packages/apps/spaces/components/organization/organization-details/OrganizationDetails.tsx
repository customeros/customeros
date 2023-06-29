import React, { ReactNode } from 'react';
import styles from './organization-details.module.scss';
import { OrganizationLocations } from '@spaces/organization/organization-locations';
import { OrganizationContacts } from '@spaces/organization/organization-contacts';
import { useOrganization } from '@spaces/hooks/useOrganization/useOrganization';
import { OrganizationBaseDetails } from '@spaces/organization/organization-details/OrganizationBaseDetails';
import { Contact } from '@spaces/graphql';

export const OrganizationDetails = ({
  id,
  children,
}: {
  id: string;
  children: ReactNode;
}) => {
  const { data, loading, error } = useOrganization({ id });

  if (error) {
    return <div>Oops! Something went wrong while loading details</div>;
  }

  return (
    <>
      <section className={styles.organizationIdCard}>
        <OrganizationBaseDetails
          id={id}
          loading={loading}
          organization={data}
        />
        <OrganizationLocations
          id={id}
          loading={loading}
          locations={data?.locations}
        />
        {children}
        <OrganizationContacts
          id={id}
          loading={loading}
          contacts={data?.contacts?.content as Array<Contact>}
        />
      </section>
    </>
  );
};

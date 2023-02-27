import React from 'react';
import styles from './organization-details.module.scss';
import { OrganizationDetailsSkeleton } from './skeletons';
import { useOrganizationDetails } from '../../../hooks/useOrganization';
import { Link } from '../../ui-kit';
export const OrganizationDetails = ({ id }: { id: string }) => {
  const { data, loading, error } = useOrganizationDetails({ id });

  if (loading) {
    return <OrganizationDetailsSkeleton />;
  }
  if (error) {
    return <>ERROR</>;
  }

  return (
    <div className={styles.organizationDetails}>
      <div className={styles.bg}>
        <div>
          <h1 className={styles.name}>{data?.name}</h1>
          <span className={styles.industry}>{data?.industry}</span>
        </div>

        <p className={styles.description}>{data?.description}</p>

        {data?.website && <Link href={data.website}> {data.website} </Link>}
      </div>
    </div>
  );
};

import React, { useState } from 'react';
import styles from './organization-details.module.scss';
import { OrganizationDetailsSkeleton } from './skeletons';
import { useOrganizationDetails } from '../../../hooks/useOrganization';
import { Button, Link } from '../../ui-kit';
import { OrganizationEdit } from './edit';
export const OrganizationDetails = ({ id }: { id: string }) => {
  const { data, loading, error } = useOrganizationDetails({ id });
  const [mode, setMode] = useState('PREVIEW');

  if (loading) {
    return <OrganizationDetailsSkeleton />;
  }
  if (error) {
    return <>ERROR</>;
  }

  if (mode === 'EDIT') {
    return <OrganizationEdit data={data} onSetMode={setMode} />;
  }

  return (
    <div className={styles.organizationDetails}>
      <div className={styles.bg}>
        <div>
          <div className={styles.header}>
            <h1 className={styles.name}>{data?.name}</h1>
            <div style={{ marginLeft: '4px' }}>
              <Button mode='secondary' onClick={() => setMode('EDIT')}>
                Edit
              </Button>
            </div>
          </div>

          <span className={styles.industry}>
            {(data?.industry ?? '')?.split('_').join(' ')}
          </span>
        </div>

        <p className={styles.description}>{data?.description}</p>

        {data?.website && <Link href={data.website}> {data.website} </Link>}
      </div>
    </div>
  );
};

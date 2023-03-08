import React, { useEffect, useState } from 'react';
import Image from 'next/image';
import { Button } from '../../ui-kit';
import styles from './contact-details.module.scss';
import { useContactPersonalDetails } from '../../../hooks/useContact';
import { ContactDetailsSkeleton } from './skeletons';
import { useRouter } from 'next/router';
import { ContactTags } from '../contact-tags';
import { useForm } from 'react-hook-form';
import { ContactPersonalDetailsEdit } from './edit';
import { getContactDisplayName } from '../../../utils';

export const ContactPersonalDetails = ({ id }: { id: string }) => {
  const router = useRouter();
  const { data, loading, error } = useContactPersonalDetails({ id });
  const [mode, setMode] = useState('PREVIEW');
  const { control, reset, setValue } = useForm();

  useEffect(() => {
    reset({
      ...data,
    });
  }, [data?.id]);

  if (loading) {
    return <ContactDetailsSkeleton />;
  }
  if (error) {
    return <>ERROR</>;
  }

  if (mode === 'EDIT') {
    return <ContactPersonalDetailsEdit data={data} onSetMode={setMode} />;
  }

  return (
    <div className={styles.header}>
      <div className={styles.photo}>
        {
          // @ts-expect-error we will have equivalent of avatar some day...
          data?.photo ? (
            // @ts-expect-error we will have equivalent of avatar some day...
            <Image src={data?.photo} alt={''} height={40} width={40} />
          ) : (
            <div>{data?.firstName?.[0]}</div>
          )
        }
      </div>
      <div className={styles.name}>
        <div className={styles.nameAndEditButton}>
          {
            //@ts-expect-error fixme later
            getContactDisplayName(data)
          }
          <div style={{ marginLeft: '4px' }}>
            <Button mode='secondary' onClick={() => setMode('EDIT')}>
              Edit
            </Button>
          </div>
        </div>

        {data?.jobRoles?.map((jobRole: any) => {
          return (
            <div
              className={styles.jobRole}
              key={`contact-job-role-${jobRole?.id}-${jobRole?.label}`}
              onClick={() =>
                router.push(`/organization/${jobRole?.organization.id}`)
              }
            >
              {jobRole?.jobTitle}

              {jobRole?.jobTitle &&
              jobRole?.organization &&
              jobRole?.organization?.name
                ? ' at'
                : ''}{' '}
              {jobRole?.organization?.name}
            </div>
          );
        })}
        <div className={styles.source}>
          <span>Source:</span>
          {data?.source || ''}
        </div>
        {/*<div className={styles.source}>*/}
        {/*  <span>Owner:</span>*/}
        {/*  {`${data?.owner?.firstName || ''} ${data?.owner?.lastName || ''}` ||*/}
        {/*    ''}*/}
        {/*</div>*/}
        <ContactTags id={id} mode='PREVIEW' />
      </div>
    </div>
  );
};

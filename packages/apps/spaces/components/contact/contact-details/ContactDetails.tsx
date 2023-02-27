import React from 'react';
import Image from 'next/image';
import { IconButton } from '../../ui-kit/atoms';
import styles from './contact-details.module.scss';
import { useContactPersonalDetails } from '../../../hooks/useContact';
import { ContactDetailsSkeleton } from './skeletons';
import { useRouter } from 'next/router';
export const ContactDetails = ({ id }: { id: string }) => {
  const router = useRouter();
  const { data, loading, error } = useContactPersonalDetails({ id });

  if (loading) {
    return <ContactDetailsSkeleton />;
  }
  if (error) {
    return <>ERROR</>;
  }

  return (
    <div className={styles.contactDetails}>
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
          <div>
            {' '}
            {data?.firstName} {data?.lastName}
          </div>
          {data?.jobRoles?.map((jobRole: any) => {
            return (
              <div
                className={styles.jobRole}
                key={`contact-job-role-${jobRole.id}-${jobRole.label}`}
                onClick={() =>
                  router.push(`/organization/${jobRole.organization.id}`)
                }
              >
                {jobRole.jobTitle}{' '}
                {jobRole.jobTitle &&
                jobRole.organization &&
                jobRole.organization.name
                  ? 'at'
                  : ''}{' '}
                {jobRole.organization.name}
              </div>
            );
          })}
          {
            <div className={styles.source}>
              <span>Source:</span>
              {data?.source || ''}
            </div>
          }
        </div>
      </div>
      <div className={styles.details}>
        <div className={styles.section}>
          <IconButton
            disabled={true}
            aria-describedby='phone-icon-label'
            mode='secondary'
            className={styles.icon}
            onClick={() => null}
            icon={
              <Image alt={''} src='/icons/phone.svg' width={20} height={20} />
            }
          />

          <div className={styles.label} id='phone-icon-label'>
            Phone
          </div>
        </div>
        <div className={styles.section}>
          <IconButton
            disabled={true}
            aria-describedby='email-icon-label'
            mode='secondary'
            className={styles.icon}
            onClick={() => null}
            icon={
              <Image
                alt={''}
                src='/icons/envelope.svg'
                width={20}
                height={20}
              />
            }
          />

          <div className={styles.label} id={'email-icon-label'}>
            Email
          </div>
        </div>
        <div className={styles.section}>
          <IconButton
            disabled={true}
            aria-describedby='message-icon-label'
            mode='secondary'
            className={styles.icon}
            onClick={() => null}
            icon={
              <Image
                alt={''}
                src='/icons/whatsapp.svg'
                width={20}
                height={20}
              />
            }
          />
          <div className={styles.label} id='message-icon-label'>
            Message
          </div>
        </div>
        <div className={styles.section}>
          <IconButton
            disabled={true}
            aria-describedby='message-icon-label'
            mode='secondary'
            className={styles.icon}
            onClick={() => null}
            icon={
              <Image
                alt={''}
                src='/icons/share-alt.svg'
                width={20}
                height={20}
              />
            }
          />
          <div className={styles.label} id='message-icon-label'>
            Share
          </div>
        </div>
      </div>
    </div>
  );
};

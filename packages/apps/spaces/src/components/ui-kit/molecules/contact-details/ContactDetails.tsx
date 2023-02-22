import React from 'react';
import Image from 'next/image';
import { Envelope, IconButton, Phone, Whatsapp } from '../../atoms';
import styles from './contact-details.module.scss';
export const ContactDetails = ({
  firstName = 'Ania',
  lastName = 'Kowalska',
  photo,
  source = 'hubspot',
  jobRoles = [
    {
      jobTitle: 'developer',
      organization: {
        name: 'Nokia',
      },
    },
    {
      jobTitle: 'developer',
      organization: {
        name: 'Orange',
      },
    },
  ],
}: any) => {
  return (
    <div className={styles.contactDetails}>
      <div className={styles.header}>
        <div className={styles.photo}>
          {photo ? (
            <Image src={photo} alt={''} height={40} width={40} />
          ) : (
            <div>
              {firstName?.[0] || ''} {lastName?.[0] || ''}
            </div>
          )}
        </div>
        <div className={styles.name}>
          <div>
            {' '}
            {firstName} {lastName}
          </div>
          {jobRoles?.map((jobRole: any) => {
            return (
              <div
                className={styles.jobRole}
                key={jobRole.id}
                onClick={
                  () => null
                  // router.push(`/organization/${jobRole.organization.id}`)
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
              {source || ''}
            </div>
          }
        </div>
      </div>
      <div className={styles.details}>
        <div className={styles.section}>
          <IconButton
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
      </div>
    </div>
  );
};

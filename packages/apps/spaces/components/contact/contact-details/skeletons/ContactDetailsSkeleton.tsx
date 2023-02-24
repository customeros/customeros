import React from 'react';
import { Skeleton } from '../../../ui-kit/atoms/skeleton';
import styles from '../contact-details.module.scss';
import Image from 'next/image';
import { IconButton } from '../../../ui-kit/atoms';

export const ContactDetailsSkeleton: React.FC = () => {
  return (
    <div className={styles.contactDetails} style={{ width: '100%' }}>
      <div className={styles.header} style={{ width: '100%' }}>
        <div className={styles.photo} style={{ background: '#ddd' }}>
          <div style={{ width: '40px', height: '40px' }} />
        </div>
        <div className={styles.name} style={{ width: '80%' }}>
          <div>
            <Skeleton height='20px' />
          </div>

          <div className={styles.jobRole}>
            <Skeleton />
          </div>

          {
            <div className={styles.source} style={{ width: '50px' }}>
              Source: <Skeleton />
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

          <Skeleton />
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

          <Skeleton />
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
          <Skeleton />
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
          <Skeleton />
        </div>
      </div>
    </div>
  );
};

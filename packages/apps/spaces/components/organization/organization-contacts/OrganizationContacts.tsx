import React from 'react';
import styles from './organization-contacts.module.scss';

import { Phone, Envelope } from '../../ui-kit';
import { OrganizationContactsSkeleton } from './skeletons';
import { useOrganizationContacts } from '../../../hooks/useOrganization';
import { useRouter } from 'next/router';
import classNames from 'classnames';

export const OrganizationContacts = ({ id }: { id: string }) => {
  const router = useRouter();
  const { data, loading, error } = useOrganizationContacts({
    id,
  });

  if (loading) {
    return <OrganizationContactsSkeleton />;
  }
  if (error) {
    return null;
  }

  return (
    <ul className={styles.contactsList}>
      {data?.map((contact) => (
        <li
          key={contact.id}
          className={classNames(styles.contactItem, styles.text)}
          role='button'
          tabIndex={0}
          onClick={() => router.push(`/contact/${contact.id}`)}
        >
          <div className={styles.personalDetails}>
            <span className={styles.name}>
              {contact?.name || `${contact?.firstName} ${contact?.lastName}`}
            </span>
            {[{ id: '234', jobTitle: 'CTO' }].map((role) => (
              <span key={role.id} className={styles.jobTitle}>
                {role.jobTitle}
              </span>
            ))}
          </div>

          {!!contact?.emails.length && (
            <div className={styles.detailsContainer}>
              <Envelope className={styles.icon} />
              {contact.emails.find((email) => email.primary)?.email ||
                contact.emails[0].email}
            </div>
          )}

          {!!contact?.phoneNumbers.length && (
            <div className={styles.detailsContainer}>
              <Phone className={styles.icon} />
              {contact.phoneNumbers.find((phoneNr) => phoneNr.primary)?.e164 ||
                contact.phoneNumbers[0].e164}
            </div>
          )}
        </li>
      ))}
    </ul>
  );
};

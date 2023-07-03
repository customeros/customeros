import React, { useRef } from 'react';
import { useRouter } from 'next/router';
import classNames from 'classnames';
import Phone from '@spaces/atoms/icons/Phone';
import Envelope from '@spaces/atoms/icons/Envelope';
import { OrganizationContactsSkeleton } from './skeletons';
import styles from './organization-contacts.module.scss';
import { ContactTags } from '@spaces/contact/contact-tags';
import { getContactDisplayName } from '../../../utils';
import { Contact } from '@spaces/graphql';

export const OrganizationContacts = ({
  loading,
  contacts,
}: {
  id: string;
  loading: boolean;
  contacts?: Array<Contact> | null;
}) => {
  const router = useRouter();
  const listRef = useRef(null);

  if (loading) {
    return <OrganizationContactsSkeleton />;
  }

  return (
    <article className={styles.contacts_section}>
      <h1 className={styles.contacts_header}>Contacts</h1>
      {!contacts?.length && (
        <div className={styles.contacts_item}>This company has no contacts</div>
      )}
      <ul className={styles.contactsList} ref={listRef}>
        {contacts?.map((contact) => (
          <li
            key={contact.id}
            className={classNames(styles.contactItem, styles.text)}
            role='button'
            tabIndex={0}
            onClick={() => router.push(`/contact/${contact.id}`)}
          >
            <div className={styles.personalDetails}>
              <span className={styles.name}>
                {getContactDisplayName(contact)}
              </span>

              {!!contact.jobRoles &&
                contact.jobRoles.map((role) => (
                  <span key={role.id} className={styles.jobTitle}>
                    {role.jobTitle}
                  </span>
                ))}
              <ContactTags
                id={contact.id}
                mode='PREVIEW'
                tags={contact?.tags}
              />
            </div>

            {!!contact?.emails.length && (
              <div className={styles.detailsContainer}>
                <Envelope className={styles.icon} width={16} height={16} />
                {contact.emails.find((email) => email.primary)?.email ||
                  contact.emails[0].email}
              </div>
            )}

            {!!contact?.phoneNumbers.length && (
              <div className={styles.detailsContainer}>
                <Phone className={styles.icon} width={16} height={16} />
                {contact.phoneNumbers.find((phoneNr) => phoneNr.primary)
                  ?.e164 || contact.phoneNumbers[0].e164}
              </div>
            )}
          </li>
        ))}
      </ul>
    </article>
  );
};

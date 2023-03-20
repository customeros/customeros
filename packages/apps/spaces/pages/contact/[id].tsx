import React from 'react';
import { DetailsPageLayout } from '../../components';
import styles from './contact.module.scss';
import { useRouter } from 'next/router';
import {
  ContactCommunicationDetails,
  ContactDetails,
  ContactEditor,
  NoteEditorModes,
} from '../../components/contact';
import { ContactPersonalDetailsCreate } from '../../components/contact/contact-details';
import ContactHistory from '../../components/contact/contact-history/ContactHistory';

function ContactDetailsPage() {
  const {
    query: { id },
    back,
    push,
  } = useRouter();

  if (id === 'new') {
    return (
      <DetailsPageLayout onNavigateBack={() => push('/')}>
        <section className={styles.idCard}>
          <ContactPersonalDetailsCreate />
        </section>
        <section className={styles.notes}></section>

        <section className={styles.timeline}></section>
      </DetailsPageLayout>
    );
  }

  return (
    <DetailsPageLayout onNavigateBack={back}>
      <section className={styles.details}>
        <ContactDetails id={id as string} />
        <ContactCommunicationDetails id={id as string} />
      </section>
      <ContactHistory id={id as string} />
      <section className={styles.notes}>
        <ContactEditor contactId={id as string} mode={NoteEditorModes.ADD} />
      </section>
    </DetailsPageLayout>
  );
}

export default ContactDetailsPage;

import React from 'react';
import { ContactDetails } from '../../components/ui-kit/molecules/contact-details';
import { ContactCommunicationDetails } from '../../components/ui-kit/molecules/contact-communication-details/ContactCommunicationDetails';
import { Button, DetailsPageLayout } from '../../components/ui-kit';
import styles from './contact.module.scss';
import { useRouter } from 'next/router';
import { ContactNoteEditor } from '../../components/contact/note-editor/NoteEditor';
function ContactDetailsPage() {
  const {
    query: { id },
  } = useRouter();

  return (
    <DetailsPageLayout>
      <section className={styles.personalDetails}>
        <ContactDetails id={id as string} />
        <ContactCommunicationDetails id={id as string} />
      </section>

      <ContactNoteEditor contactId={id as string} />
      <section className={styles.timeline}></section>
    </DetailsPageLayout>
  );
}

export default ContactDetailsPage;

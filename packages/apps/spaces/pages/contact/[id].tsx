import React from 'react';
import { DetailsPageLayout } from '../../components/ui-kit';
import styles from './contact.module.scss';
import { useRouter } from 'next/router';
import {
  ContactCommunicationDetails,
  ContactDetails,
  ContactNoteEditor,
  NoteEditorModes,
} from '../../components/contact';

function ContactDetailsPage() {
  const {
    query: { id },
    back,
  } = useRouter();

  return (
    <DetailsPageLayout onNavigateBack={back}>
      <section className={styles.personalDetails}>
        <ContactDetails id={id as string} />
        <ContactCommunicationDetails id={id as string} />
      </section>
      <section className={styles.notes}>
        <ContactNoteEditor
          contactId={id as string}
          mode={NoteEditorModes.ADD}
        />
      </section>

      <section className={styles.timeline}></section>
    </DetailsPageLayout>
  );
}

export default ContactDetailsPage;

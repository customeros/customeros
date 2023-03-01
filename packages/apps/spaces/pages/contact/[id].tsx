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
import { ContactPersonalDetailsCreate } from '../../components/contact/contact-details';
import ContactHistory from "../../components/contact/contact-history/ContactHistory";

function ContactDetailsPage() {
  const {
    query: { id },
    back,
  } = useRouter();

  if (id === 'new') {
    return (
      <DetailsPageLayout onNavigateBack={back}>
        <section className={styles.personalDetails}>
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
            <section className={styles.notes}>
                <ContactNoteEditor
                    contactId={id as string}
                    mode={NoteEditorModes.ADD}
                />
            </section>
            <section className={styles.timeline}>
                <ContactHistory id={id as string} />
            </section>
    </DetailsPageLayout>
  );
}

export default ContactDetailsPage;

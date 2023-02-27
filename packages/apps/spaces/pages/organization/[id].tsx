import React from 'react';
import { DetailsPageLayout } from '../../components/ui-kit';
import styles from './organization.module.scss';
import { useRouter } from 'next/router';
import {
  OrganizationDetails,
  OrganizationNoteEditor,
  NoteEditorModes,
  OrganizationContacts,
} from '../../components/organization';

function OrganizationDetailsPage() {
  const {
    query: { id },
    push,
  } = useRouter();

  return (
    <DetailsPageLayout onNavigateBack={() => push('/')}>
      <section className={styles.organizationDetails}>
        <OrganizationDetails id={id as string} />
        <OrganizationContacts id={id as string} />
      </section>
      <section className={styles.notes}>
        <OrganizationNoteEditor
          organizationId={id as string}
          mode={NoteEditorModes.ADD}
        />
      </section>
      <section className={styles.timeline}></section>
    </DetailsPageLayout>
  );
}

export default OrganizationDetailsPage;

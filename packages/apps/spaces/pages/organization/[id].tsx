import React from 'react';
import { DetailsPageLayout } from '../../components/ui-kit';
import styles from './organization.module.scss';
import { useRouter } from 'next/router';
import {
  OrganizationDetails,
  OrganizationEditor,
  NoteEditorModes,
  OrganizationContacts,
  OrganizationCreate,
} from '../../components/organization';
import { OrganizationTimeline } from '../../components/organization/organization-timeline';

function OrganizationDetailsPage() {
  const {
    query: { id },
    push,
  } = useRouter();

  if (id === 'new') {
    return (
      <DetailsPageLayout onNavigateBack={() => push('/')}>
        <section className={styles.organizationIdCard}>
          <OrganizationCreate />
        </section>
        <section className={styles.notes}></section>
        <section className={styles.timeline}></section>
      </DetailsPageLayout>
    );
  }
  return (
    <DetailsPageLayout onNavigateBack={() => push('/')}>
      <section className={styles.organizationIdCard}>
        <OrganizationDetails id={id as string} />
      </section>
      <section className={styles.organizationDetails}>
        <OrganizationContacts id={id as string} />
      </section>
      <section className={styles.notes}>
        <OrganizationEditor
          organizationId={id as string}
          mode={NoteEditorModes.ADD}
        />
      </section>
      <section className={styles.timeline}>
        <OrganizationTimeline id={id as string} />
      </section>
    </DetailsPageLayout>
  );
}

export default OrganizationDetailsPage;

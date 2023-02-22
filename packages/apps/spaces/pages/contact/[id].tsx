import React from 'react';
import { ContactDetails } from '../../src/components/ui-kit/molecules/contact-details';
import { ContactCommunicationDetails } from '../../src/components/ui-kit/molecules/contact-communication-details/ContactCommunicationDetails';
import { DetailsPageLayout } from '../../src/components/ui-kit';
import styles from './contact.module.scss';
function ContactDetailsPage() {
  return (
    <DetailsPageLayout>
      <section className={styles.personalDetails}>
        <ContactDetails
          name='Details'
          photo=''
          phone='40294-02349-234'
          email='asdsad@op.pl'
          address={''}
          birthday={'11/1/1996'}
          socialProfiles={[]}
          notes={[]}
        />
        <ContactCommunicationDetails />
      </section>
      <section className={styles.notes}>placeholder</section>
      <section className={styles.timeline}>placeholder</section>
    </DetailsPageLayout>
  );
}

export default ContactDetailsPage;

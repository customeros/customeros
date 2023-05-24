import React from 'react';
import { DetailsPageLayout } from '@spaces/layouts/details-page-layout';
import { motion } from 'framer-motion';
import styles from '../../../pages/organization/organization.module.scss';
import { TimelineSkeleton } from '@spaces/organisms/timeline';
import { ContactDetailsSkeleton } from '@spaces/contact/contact-details/skeletons';
import { ContactCommunicationDetailsSkeleton } from '@spaces/contact/contact-communication-details/skeletons';

const easing = [0.6, -0.05, 0.01, 0.99];

const fadeInUp = {
  initial: {
    opacity: 0,
    transition: { duration: 0.6, ease: easing },
  },
  animate: {
    opacity: 1,
    transition: {
      duration: 0.6,
      ease: easing,
    },
  },
};
export const ContactProfileSkeleton: React.FC = () => {
  return (
    <DetailsPageLayout>
      <motion.section variants={fadeInUp}>
        <ContactDetailsSkeleton />
        <ContactCommunicationDetailsSkeleton />
      </motion.section>
      <motion.section variants={fadeInUp}>
        <TimelineSkeleton />
      </motion.section>
      <motion.section
        variants={fadeInUp}
        className={styles.notes}
      ></motion.section>
    </DetailsPageLayout>
  );
};

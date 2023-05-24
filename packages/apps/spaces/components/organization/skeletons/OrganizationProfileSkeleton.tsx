import React from 'react';
import { DetailsPageLayout } from '@spaces/layouts/details-page-layout';
import { motion } from 'framer-motion';
import styles from '../../../pages/organization/organization.module.scss';
import { OrganizationDetailsSkeleton } from '@spaces/organization/organization-details/skeletons';
import { OrganizationContactsSkeleton } from '@spaces/organization/organization-contacts/skeletons';
import { TimelineSkeleton } from '@spaces/organisms/timeline/skeletons/TimelineSkeleton';

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
export const OrganizationProfileSkeleton: React.FC = () => {
  return (
    <DetailsPageLayout>
      <motion.section variants={fadeInUp} className={styles.organizationIdCard}>
        <OrganizationDetailsSkeleton />
      </motion.section>
      <motion.section
        variants={fadeInUp}
        className={styles.organizationDetails}
      >
        <OrganizationContactsSkeleton />
      </motion.section>
      <motion.section
        variants={fadeInUp}
        className={styles.notes}
      ></motion.section>
      <section className={styles.timeline}>
        <TimelineSkeleton />
      </section>
    </DetailsPageLayout>
  );
};

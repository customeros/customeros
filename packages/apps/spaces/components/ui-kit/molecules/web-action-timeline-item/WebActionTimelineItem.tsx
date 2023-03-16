import React from 'react';
import format from 'date-fns/format';
import Link from 'next/link';

import styles from './web-action-timeline-item.module.scss';
import { Globe } from '../../atoms';
import { capitalizeFirstLetter, DateTimeUtils } from '../../../../utils';

// interface Props extends ContactWebAction {
//     contactName?: string
// }

export const WebActionTimelineItem = ({
  pageTitle,
  pageUrl,
  engagedTime,
  contactName,
}: any): JSX.Element => {
  return (
    <article className={`${styles.actionContainer}`}>
      {contactName && <span className={styles.visitor}>{contactName}</span>}
      <span> visited: </span>
      <Link
        href={pageUrl}
        target={'_blank'}
        style={{ textDecoration: 'underline' }}
      >
        {pageTitle && pageTitle !== '' && <span>{pageTitle}</span>}
        {(!pageTitle || pageTitle === '') && <span>{pageUrl}</span>}
      </Link>

      <span> - {DateTimeUtils.formatSecondsDuration(engagedTime)}</span>
    </article>
  );
};

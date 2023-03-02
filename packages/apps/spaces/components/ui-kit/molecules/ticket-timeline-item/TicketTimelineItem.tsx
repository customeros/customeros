import React from 'react';

import styles from './ticket-timeline-item.module.scss';
import Ticket from '../../atoms/icons/Ticket';
import { TagsList } from '../../atoms';

// interface Props extends ContactWebAction {
//     contactName?: string
// }

export const TicketTimelineItem = ({
  createdAt,
  updatedAt,
  subject,
  status,
  priority,
  description,
  tags,
  ...rest
}: any): JSX.Element => {
  return (
    <div className={styles.x}>
      <Ticket />
      <article className={`${styles.ticketContainer}`}>
        <div className={`${styles.ticketHeader}`}>
          <div className={`${styles.ticketHeaderSubject}`}>{subject}</div>
          <div className={`${styles.ticketHeaderStatus}`}>{status}</div>
        </div>

        {tags && tags.length > 0 && (
          <div className={`${styles.tags}`}>
            <TagsList tags={tags ?? []} readOnly={true} />
          </div>
        )}

        <div>{description}</div>
      </article>
    </div>
  );
};

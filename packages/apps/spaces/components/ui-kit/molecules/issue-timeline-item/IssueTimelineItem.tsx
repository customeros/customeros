import React from 'react';
import styles from './issue-timeline-item.module.scss';
import Ticket from '@spaces/atoms/icons/Ticket';
import { TagsList } from '@spaces/atoms/tags/TagList';
import { Issue } from '@spaces/graphql';

export const IssueTimelineItem = ({
  tags,
  status,
  subject,
  description,
}: Issue): JSX.Element => {
  const parsedTags =
    tags
      ?.filter(Boolean)
      .map((tag) => ({ id: tag?.id ?? '', name: tag?.name ?? '' })) ?? [];

  return (
    <div className={styles.x}>
      <article className={`${styles.ticketContainer}`}>
        <div className={`${styles.ticketHeader}`}>
          <div className={`${styles.ticketHeaderSubject}`}>
            <Ticket className={`${styles.ticketHeaderPicture}`} />
            {subject}
          </div>
          <div className={`${styles.ticketHeaderStatus}`}>{status}</div>
        </div>

        {tags && tags.length > 0 && (
          <div className={`${styles.tags}`}>
            <TagsList tags={parsedTags} readOnly={true} />
          </div>
        )}

        <div>{description}</div>
      </article>
    </div>
  );
};

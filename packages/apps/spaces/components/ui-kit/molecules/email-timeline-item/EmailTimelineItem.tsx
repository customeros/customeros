import React, { ReactNode, useRef, useState } from 'react';
import sanitizeHtml from 'sanitize-html';
import styles from './email-timeline-item.module.scss';
import { Button } from '../../atoms';
import linkifyHtml from 'linkify-html';

interface Props {
  emailContent: string;
  emailContentType: string;
  sender: string;
  recipients: string | Array<string>;
  cc?: string | Array<string>;
  bcc?: string | Array<string>;
  subject: string;
  children?: ReactNode;
}

export const EmailTimelineItem: React.FC<Props> = ({
  emailContent,
  emailContentType,
  sender,
  recipients,
  subject,
  cc,
  bcc,
  children,
}) => {
  const [expanded, toggleExpanded] = useState(false);
  const timelineItemRef = useRef<HTMLDivElement>(null);

  const handleToggleExpanded = () => {
    toggleExpanded(!expanded);
    if (timelineItemRef?.current && expanded) {
      timelineItemRef?.current?.scrollIntoView();
    }
  };
  return (
    <div className={styles.emailWrapper}>
      <div ref={timelineItemRef} className={styles.scrollToView} />
      <article className={`${styles.emailContainer}`}>
        <div className={styles.emailData}>
          <table className={styles.emailDataTable}>
            <tr>
              <th className={styles.emailParty}>From:</th>
              <td>{sender}</td>
            </tr>
            <tr>
              <th className={styles.emailParty}>To:</th>
              <td>
                {
                  <div className={styles.emailRecipients}>
                    {typeof recipients === 'string'
                      ? recipients
                      : recipients.map((recipient) => (
                          <span
                            className={styles.emailRecipient}
                            key={recipient}
                          >
                            {recipient}
                          </span>
                        ))}
                  </div>
                }
              </td>
            </tr>

            {!!cc?.length && (
              <tr>
                <th className={styles.emailParty}>CC:</th>
                <td>
                  {
                    <div className={styles.emailRecipients}>
                      {typeof cc === 'string'
                        ? cc
                        : cc.map((recipient) => (
                            <span
                              className={styles.emailRecipient}
                              key={recipient}
                            >
                              {recipient}
                            </span>
                          ))}
                    </div>
                  }
                </td>
              </tr>
            )}
            {!!bcc?.length && (
              <tr>
                <th className={styles.emailParty}>BCC:</th>
                <td>
                  {
                    <div className={styles.emailRecipients}>
                      {typeof bcc === 'string'
                        ? bcc
                        : bcc.map((recipient) => (
                            <span
                              className={styles.emailRecipient}
                              key={recipient}
                            >
                              {recipient}
                            </span>
                          ))}
                    </div>
                  }
                </td>
              </tr>
            )}
            <tr>
              <th className={styles.emailParty}>Subject:</th>
              <td>{subject}</td>
            </tr>
          </table>

          <div className={styles.stamp}>
            <div />
          </div>
        </div>
        <div
          className={`${styles.emailContentContainer} ${
            !expanded ? styles.eclipse : ''
          }`}
          style={{ height: expanded ? '100%' : '80px' }}
        >
          {emailContentType === 'text/html' && (
            <div
              className={`text-overflow-ellipsis ${styles.emailContent}`}
              dangerouslySetInnerHTML={{
                __html: sanitizeHtml(
                    linkifyHtml(emailContent, {
                      defaultProtocol: 'https',
                      rel: 'noopener noreferrer',
                    }),
                ),
              }}
            ></div>
          )}
          {emailContentType === 'text/plain' && (
            <div className={`text-overflow-ellipsis ${styles.emailContent}`}>
              {emailContent}
            </div>
          )}

          {!expanded && <div className={styles.eclipse} />}
        </div>
        <div className={styles.toggleExpandButton}>
          <Button onClick={() => handleToggleExpanded()} mode='link'>
            {expanded ? 'Collapse' : 'Expand'}
          </Button>
          {children}
        </div>
      </article>
    </div>
  );
};

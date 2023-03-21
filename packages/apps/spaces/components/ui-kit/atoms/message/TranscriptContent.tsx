import * as React from 'react';
import styles from './message.module.scss';
import linkifyHtml from 'linkify-html';
import { ReactNode } from 'react';
import classNames from 'classnames';

interface Content {
  type?: string;
  mimetype: string;
  body: string;
}
interface TranscriptElement {
  party: any;
  text: string;
}

interface TranscriptContentProps {
  response: Array<Content>;
  children?: ReactNode;
}

export const TranscriptContent: React.FC<TranscriptContentProps> = ({
  response,
  children,
}) => {
  return (
    <div className={styles.transcriptionContainer}>
      {response?.map(
        // @ts-expect-error fixme
        (transcriptElement: TranscriptElement, index: number) => {
          const textWithLinks = linkifyHtml(transcriptElement.text, {
            defaultProtocol: 'https',
            rel: 'noopener noreferrer',
          });
          return (
            <div
              key={index}
              className={classNames(styles.singleMessage, {
                [styles.isleft]: transcriptElement.party.tel,
                [styles.isright]: !transcriptElement.party.tel,
              })}
            >
              <div
                className={classNames(styles.channelIcon, {
                  [styles.channelIconShown]: index === 0,
                })}
              >
                {index === 0 && children && children}
              </div>

              <div
                className={classNames(styles.message, {
                  [styles.left]: transcriptElement.party.tel,
                  [styles.right]: !transcriptElement.party.tel,
                })}
                style={{ width: '60%' }}
              >
                {textWithLinks}
              </div>
            </div>
          );
        },
      )}
    </div>
  );
};

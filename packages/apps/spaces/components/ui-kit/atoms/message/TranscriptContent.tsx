import * as React from 'react';
import styles from './message.module.scss';
import linkifyHtml from 'linkify-html';

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
}

export const TranscriptContent: React.FC<TranscriptContentProps> = ({
  response,
}) => {
  const reversed = [...response].reverse();
  return (
    <div className={styles.transcriptionContainer}>
      {reversed?.map(
        // @ts-expect-error fixme
        (transcriptElement: TranscriptElement, index: number) => {
          const textWithLinks = linkifyHtml(transcriptElement.text, {
            defaultProtocol: 'https',
            rel: 'noopener noreferrer',
          });
          return (
            <div key={index}>
              {transcriptElement.party.tel && (
                <div
                  className={`${styles.message} ${styles.left} ${styles.test}`}
                >
                  {textWithLinks}
                </div>
              )}
              {!transcriptElement.party.tel && (
                <div
                  className={`${styles.message} ${styles.right}`}
                  style={{ background: '#C5EDCE', borderRadius: '5px' }}
                >
                  {textWithLinks}
                </div>
              )}
            </div>
          );
        },
      )}
    </div>
  );
};

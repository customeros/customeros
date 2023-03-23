import * as React from 'react';
import styles from './message.module.scss';
import linkifyHtml from 'linkify-html';
import { ReactNode } from 'react';
import classNames from 'classnames';

interface TranscriptElement {
  party: any;
  text: string;
}

interface TranscriptContentProps {
  messages: Array<TranscriptElement>;
  children?: ReactNode;
  firstIndex: {
    received: number;
    send: number;
  };
}

export const TranscriptContent: React.FC<TranscriptContentProps> = ({
  messages = [],
  children,
  firstIndex,
}) => {
  return (
    <>
      {messages?.map((transcriptElement: TranscriptElement, index: number) => {
        const showIcon =
          index === firstIndex.send || index === firstIndex.received;
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
                [styles.channelIconShown]: showIcon,
              })}
            >
              {showIcon && children}
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
      })}
    </>
  );
};

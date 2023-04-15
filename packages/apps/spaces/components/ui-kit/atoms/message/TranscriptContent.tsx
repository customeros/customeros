import * as React from 'react';
import styles from './message.module.scss';
import linkifyHtml from 'linkify-html';
import linkifyStr from 'linkify-string';
import { ReactNode } from 'react';
import classNames from 'classnames';
import sanitizeHtml from 'sanitize-html';
import AudioPlayer from 'react-audio-player';

interface TranscriptElement {
  party: any;
  text: string;
  file_id?: string;
}

interface TranscriptContentProps {
  messages: Array<TranscriptElement>;
  children?: ReactNode;
  firstIndex: {
    received: number | null;
    send: number | null;
  };
  contentType?: string;
}

export const TranscriptContent: React.FC<TranscriptContentProps> = ({
  messages = [],
  children,
  firstIndex,
  contentType,
}) => {
  return (
    <>
      {messages?.map((transcriptElement: TranscriptElement, index: number) => {
        const showIcon =
          index === firstIndex.send || index === firstIndex.received;
        const transcriptContent =
          transcriptElement?.text && contentType === 'text/html'
            ? transcriptElement.text
            : `<p>${transcriptElement.text}</p>`;

        return (
          <div
            key={index}
            className={classNames(styles.singleMessage, {
              [styles.isleft]: transcriptElement?.party.tel,
              [styles.isright]: !transcriptElement?.party.tel,
            })}
          >
            <div
              className={classNames(styles.channelIcon, {
                [styles.channelIconShown]: showIcon,
              })}
            >
              {showIcon && children}
            </div>

            {transcriptElement?.text && (
              <div
                className={classNames(styles.message, {
                  [styles.left]: transcriptElement?.party.tel,
                  [styles.right]: !transcriptElement?.party.tel,
                })}
                style={{ width: '60%' }}
                dangerouslySetInnerHTML={{
                  __html: sanitizeHtml(
                    linkifyHtml(transcriptContent, {
                      defaultProtocol: 'https',
                      rel: 'noopener noreferrer',
                    }),
                    {
                      allowedTags: sanitizeHtml.defaults.allowedTags.concat([
                        'img',
                      ]),
                      allowedAttributes: {
                        img: ['src', 'alt'],
                        a: ['href', 'rel'],
                      },
                      allowedSchemes: ['data', 'http', 'https'],
                    },
                  ),
                }}
              ></div>
            )}
            {transcriptElement?.file_id && (
              <AudioPlayer
              src={"/fs/file/" + transcriptElement.file_id + "/download?inline=true"} controls autoPlay/>
            )}
          </div>
        );
      })}
    </>
  );
};

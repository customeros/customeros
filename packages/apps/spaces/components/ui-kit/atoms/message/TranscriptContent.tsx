import * as React from 'react';
import styles from './message.module.scss';
import linkifyHtml from 'linkify-html';
import linkifyStr from 'linkify-string';
import { ReactNode } from 'react';
import classNames from 'classnames';
import sanitizeHtml from 'sanitize-html';
import {
  IconButton,
  Play,
} from '../../atoms';

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

  const playerButtonActive = new Map<number, [boolean, React.Dispatch<React.SetStateAction<boolean>>]>();
  messages?.map((transcriptElement: TranscriptElement, index: number) => {
      const [state, setState] = React.useState(false);
      playerButtonActive.set(index, [state, setState])

  });

  return (
    <>
      {messages?.map((transcriptElement: TranscriptElement, index: number) => {
        const buttonState = playerButtonActive.get(index)
        const [showPlayer, setShowPlayer] = buttonState ? buttonState : [false, () => {}]
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
            {!transcriptElement?.text && transcriptElement?.file_id && (
              <div
                className={classNames(styles.message, {
                  [styles.left]: transcriptElement?.party.tel,
                  [styles.right]: !transcriptElement?.party.tel,
                })}>
                  <i>*Unable to Transcribe Audio*</i>
                </div>
            )}
            {transcriptElement?.file_id && (
              <IconButton onClick={() => setShowPlayer(!showPlayer)}
                icon={<Play/>}
                style={{ marginBottom: 0, color: 'green' }}
              />

            )}
            {transcriptElement?.file_id && showPlayer && (
              <audio src={"/fs/file/" + transcriptElement.file_id + "/download?inline=true"} autoPlay/>
            )}
          </div>
        );
      })}
    </>
  );
};

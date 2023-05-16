import * as React from 'react';
import styles from './message.module.scss';
import {VolumeUp} from '../icons';
import classNames from 'classnames';
import sanitizeHtml from 'sanitize-html';
import linkifyHtml from 'linkify-html';
import { IconButton } from '../icon-button';
import { ReactNode, useState } from 'react';
interface TranscriptElement {
  party: any;
  text: string;
  file_id?: string;
}

interface Props {
  transcriptElement: TranscriptElement;
  index: number;
  children?: ReactNode;
  firstIndex?: {
    received: number | null;
    send: number | null;
  };
  contentType?: string;
  isLeft: boolean;
  showAvatar?: boolean;
}

export const Message = ({
  transcriptElement,
  index,
  contentType,
  firstIndex,
  children,
  isLeft,
  showAvatar,
}: Props) => {
  const [showPlayer, setShowPlayer] = useState(false);

  const showIcon =
    !firstIndex || index === firstIndex?.send || index === firstIndex?.received;
  const transcriptContent =
    transcriptElement?.text && contentType === 'text/html'
      ? transcriptElement.text
      : `<p>${transcriptElement.text}</p>`;
  return (
    <div
      key={index}
      className={classNames(styles.singleMessage, {
        [styles.isleft]: isLeft,
        [styles.isright]: !isLeft,
      })}
    >
      <div
        className={classNames({
          [styles.channelIcon]: !showAvatar,
          [styles.channelIconShown]: showIcon && !showAvatar,
        })}
      >
        {showIcon && children}
      </div>

      {transcriptElement?.text && (
        <div
          className={classNames(styles.message, {
            [styles.left]: isLeft,
            [styles.right]: !isLeft,
          })}
          style={{ width: '60%' }}
          dangerouslySetInnerHTML={{
            __html: sanitizeHtml(
              linkifyHtml(transcriptContent, {
                defaultProtocol: 'https',
                rel: 'noopener noreferrer',
              }),
              {
                allowedTags: sanitizeHtml.defaults.allowedTags.concat(['img']),
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
            [styles.left]: isLeft,
            [styles.right]: !isLeft,
          })}
        >
          <i>*Unable to Transcribe Audio*</i>
        </div>
      )}
      {transcriptElement?.file_id && (
        <IconButton
          mode={'text'}
          label='Play'
          onClick={() => setShowPlayer(!showPlayer)}
          icon={<VolumeUp />}
          style={{ marginBottom: 0, color: 'green', width: 40, height: 35 }}
        />
      )}
      {transcriptElement?.file_id && showPlayer && (
        <audio
          src={
            '/fs/file/' + transcriptElement.file_id + '/download?inline=true'
          }
          autoPlay
        />
      )}
    </div>
  );
};

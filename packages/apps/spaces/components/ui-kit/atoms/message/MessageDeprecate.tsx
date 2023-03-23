import * as React from 'react';
import styles from './message.module.scss';
import { DialogContent } from './DialogContent';
import { AnalysisContent } from './AnalysisContent';
import { MessageIcon, Phone } from '../icons';
import classNames from 'classnames';

interface Props {
  message: any;

  date: any;
  index: number;
  mode: 'CHAT' | 'PHONE_CALL' | 'LIVE_CONVERSATION';
}

export const MessageDeprecate = ({ message, mode, index }: Props) => {
  const decodeContent = (content: string) => {
    let response;
    try {
      response = JSON.parse(content);
    } catch (e) {
      response = {
        dialog: {
          type: 'MESSAGE',
          mimetype: 'text/plain',
          body: content,
        },
      };
    }
    return response;
  };

  const content = decodeContent(message.content);

  if (content.analysis && message.direction) {
    message.direction = 3;
  }

  return (
    <div style={{ width: '100%' }}>
      {content.dialog && message.direction === 0 && (
        <div className={styles.dialogMessageLeft}>
          <div
            className={styles.messageContainer}
            style={{
              marginBottom: 'var(--spacing-xxs)',
              justifyContent: 'flex-start',
            }}
          >
            <div
              className={`${styles.message} ${styles.left}`}
              style={{ width: '60%' }}
            >
              <div className='flex'>
                <div
                  className={classNames(
                    styles.channelIcon,
                    styles.channelIconLeft,
                    {
                      [styles.channelIconShown]: index === 0,
                      [styles.chatIcon]: mode === 'CHAT',
                      [styles.phoneCall]: mode === 'PHONE_CALL',
                    },
                  )}
                >
                  {/*{console.log(feedId, index, mode, content.dialog.body)}*/}
                  {mode === 'CHAT' && index === 0 && <MessageIcon />}
                  {mode === 'PHONE_CALL' && index === 0 && <Phone />}
                </div>
                <DialogContent dialog={content.dialog} />
              </div>
            </div>
          </div>
        </div>
      )}
      {content.dialog && message.direction === 1 && (
        <div className={styles.dialogMessageRight}>
          <div
            className={styles.messageContainer}
            style={{
              marginBottom: 'var(--spacing-xxs)',
              justifyContent: 'flex-end',
            }}
          >
            <div className={`${styles.message} ${styles.right}`}>
              <DialogContent dialog={content.dialog} />
            </div>
            <div
              className={classNames(
                styles.channelIcon,
                styles.channelIconRight,
                {
                  [styles.channelIconShown]:
                    mode === 'PHONE_CALL' ? index === 1 : index === 0,
                  [styles.chatIcon]: mode === 'CHAT',
                  [styles.phoneCall]: mode === 'PHONE_CALL',
                },
              )}
            >
              {mode === 'CHAT' && index === 0 && <MessageIcon />}
              {mode === 'PHONE_CALL' && index === 1 && <Phone />}
            </div>
          </div>
        </div>
      )}
      {content.analysis?.type !== 'summary' && (
        <>
          <div
            className={`${styles.center}`}
            style={{ background: 'transparent', borderRadius: '5px' }}
          >
            <AnalysisContent analysis={content.analysis}>
              {mode === 'CHAT' ? <MessageIcon /> : <Phone />}
            </AnalysisContent>
          </div>
        </>
      )}
    </div>
  );
};

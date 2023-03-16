import * as React from 'react';
import styles from './message.module.scss';
import { DialogContent } from './DialogContent';
import { AnalysisContent } from './AnalysisContent';

// import { ConversationItem } from '../../../models/conversation-item';
interface Props {
  message: any;
  sender: {
    loaded: boolean;
    email: string;
    firstName: string;
    lastName: string;
    phoneNumber: string;
  };
  date: any;
  previousMessage: number | null;
  index: number;
}

export const Message = ({
  message,
  sender,
  date,
  previousMessage,
  index,
}: Props) => {
  const decodeChannel = () => {
    return 'Voice call';
  };

  interface Content {
    type?: string;
    mimetype: string;
    body: string;
  }

  interface MiniVcon {
    parties: Array<VConParty>;
    dialog?: Content;
    analysis?: Content;
  }

  interface VConParty {
    tel?: string;
    stir?: string;
    mailto?: string;
    name?: string;
  }
  const getUser = (msg: MiniVcon): VConParty => {
    if (msg.parties) {
      for (const party of msg.parties) {
        if (party.mailto) {
          return party;
        }
      }
    }
    return { mailto: 'unknown' };
  };

  const getContact = (msg: MiniVcon): VConParty => {
    console.log('getContact', msg);
    if (msg.parties) {
      for (const party of msg.parties) {
        if (party.tel) {
          return party;
        }
      }
    }
    return { tel: 'unknown' };
  };

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
        <>
          {(index === 0 ||
            (index > 0 && previousMessage !== message?.direction)) && (
            <div className='mb-1 text-gray-600 part-1'>
              {sender.firstName && sender.lastName && (
                <>
                  {sender.firstName} {sender.lastName}
                </>
              )}
              {!sender.firstName && !sender.lastName && <>{sender.email}</>}
              {!sender.firstName && !sender.lastName && sender.phoneNumber && (
                <>{sender.phoneNumber}</>
              )}
              HELLO
            </div>
          )}
          <div className={styles.messageContainer}>
            <div className={`${styles.message} ${styles.left}`}>
              <DialogContent dialog={content.dialog} />

              <div
                className='flex align-content-end'
                style={{
                  width: '100%',
                  textAlign: 'right',
                  fontSize: '12px',
                  color: '#C1C1C1',
                }}
              >
                <span className='flex-grow-1'></span>
              </div>
            </div>
          </div>
        </>
      )}
      {content.dialog && message.direction === 1 && (
        <>
          {(index === 0 ||
            (index > 0 && previousMessage !== message?.direction)) && (
            <div className='w-full flex'>
              <div className='flex-grow-1'></div>
              {
                <div className='flex-grow-0 mb-1 pr-3'>
                  {sender.firstName && sender.lastName && (
                    <>
                      {sender.firstName} {sender.lastName}
                    </>
                  )}
                  {!sender.firstName && !sender.lastName && sender.email && (
                    <>{sender.email}</>
                  )}
                  {!sender.firstName &&
                    !sender.lastName &&
                    sender.phoneNumber && <>{sender.phoneNumber}</>}
                </div>
              }
            </div>
          )}

          <div
            className={styles.messageContainer}
            style={{ justifyContent: 'flex-end' }}
          >
            <div
              className={`${styles.message} ${styles.right}`}
              // style={{ background: '#C5EDCE', borderRadius: '5px' }}
            >
              <DialogContent dialog={content.dialog} />
              <div
                className='flex align-content-end'
                style={{
                  width: '100%',
                  textAlign: 'right',
                  fontSize: '12px',
                  color: '#C1C1C1',
                }}
              >
                <span className='flex-grow-1'></span>
                <span className='text-gray-600 mr-2'>{decodeChannel()}</span>
              </div>
            </div>
          </div>
        </>
      )}
      {content.analysis && (
        <>
          <div className='w-full flex'>
            <div className='flex-grow-1'></div>
            {
              <div className='flex flex-row flex-grow-0'>
                <div className={'mr-1'}>{getContact(content).tel}</div>
                <div className={'mr-1'} style={{ textAlign: 'center' }}>
                  <b>{content.analysis.type}</b>
                </div>
                <div className={'mr-1'}>{getUser(content).mailto}</div>
              </div>
            }
          </div>

          <div
            className={`${styles.message} ${styles.center}`}
            style={{ background: 'transparent', borderRadius: '5px' }}
          >
            <AnalysisContent analysis={content.analysis} />
            <div
              className='flex align-content-end'
              style={{
                width: '100%',
                textAlign: 'center',
                fontSize: '12px',
                color: '#C1C1C1',
              }}
            >
              <span className='flex-grow-1'></span>
              <span className='text-gray-600 mr-2'>{decodeChannel()}</span>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

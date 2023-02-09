import * as React from 'react';
import styles from './message.module.css';
import { EmailTimelineItem } from '../../molecules';
// import { ConversationItem } from '../../../models/conversation-item';
interface Props {
  message: any;
  feedInitiator: {
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
  feedInitiator,
  date,
  previousMessage,
  index,
}: Props) => {
  const decodeChannel = (channel: number) => {
    switch (channel) {
      case 0:
        return 'Web chat';
      case 1:
        return 'Email';
      case 2:
        return 'WhatsApp';
      case 3:
        return 'Facebook';
      case 4:
        return 'Twitter';
      case 5:
        return 'Phone call';
    }
    return '';
  };

  return (
    <>
      {message.direction == 0 && (
        <>
          {index === 0 && (
            <div className='mb-1 text-gray-600'>
              {feedInitiator.firstName && feedInitiator.lastName && (
                <>
                  {feedInitiator.firstName} {feedInitiator.lastName}
                </>
              )}
              {!feedInitiator.firstName && !feedInitiator.lastName && (
                <>{feedInitiator.email}</>
              )}
            </div>
          )}
          <div className={styles.messageContainer}>
            <div className={`${styles.message} ${styles.left}`}>
              {message.content}
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
                <span className='text-gray-600 mr-2'>
                  {decodeChannel(message.type)}
                </span>
                {/*{date && (*/}
                {/*  // <Moment*/}
                {/*  //   className='text-sm text-gray-600'*/}
                {/*  //   date={date}*/}
                {/*  //   format={'HH:mm'}*/}
                {/*  // ></Moment>*/}
                {/*)}*/}
              </div>
            </div>
          </div>
        </>
      )}
      {message.direction == 1 && (
        <>
          {(index === 0 ||
            (index > 0 && previousMessage !== message?.direction)) && (
            <div className='w-full flex'>
              <div className='flex-grow-1'></div>
              {/*<div className="flex-grow-0 mb-1 pr-3">To be added</div>*/}
            </div>
          )}

          <div
            className={styles.messageContainer}
            style={{ justifyContent: 'flex-end' }}
          >
            <div
              className={`${styles.message} ${styles.right}`}
              style={{ background: '#C5EDCE', borderRadius: '5px' }}
            >
              {message.content}
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
                <span className='text-gray-600 mr-2'>
                  {decodeChannel(message.type)}
                </span>
                {/*{date && (*/}
                {/*  // <Moment*/}
                {/*  //   className='text-sm text-gray-600'*/}
                {/*  //   date={date}*/}
                {/*  //   format={'HH:mm'}*/}
                {/*  // ></Moment>*/}
                {/*)}*/}
              </div>
            </div>
          </div>
        </>
      )}
    </>
  );
};

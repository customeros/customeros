import * as React from 'react';
import styles from './message.module.css';
import sanitizeHtml from 'sanitize-html';
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

  interface Content {
    type?: string;
    mimetype: string;
    body: string;
  }
  interface MiniVcon {
    dialog?: Content
    analysis?: Content
  }

  interface TranscriptElement {
    party: string;
    text: string;
  }

  const decodeContent = (content: string)  => {
    let response;
    try {
      response = JSON.parse(content);
      console.log("***************Managed to parse JSON!!!!" + content)
      console.log("Got result:" + JSON.stringify(response));
    } catch (e) {
      response = {dialog: 
       {
          type: 'MESSAGE',
          mimetype: 'text/plain',
          body: content
       } };
    }
    //console.log("Got result:" + JSON.stringify(response));
    return response;
  }

  const displayContent = (content: MiniVcon) => {
    if (content.dialog) {
      if (content.dialog.mimetype === 'text/plain') {
        return content.dialog.body
      } else if (content.dialog.mimetype === 'text/html') {
        return <div
        className={`text-overflow-ellipsis ${styles.emailContent}`}
        dangerouslySetInnerHTML={{ __html: sanitizeHtml(content.dialog.body) }}
      ></div>
      }
    }
    if (content.analysis) {
      if (content.analysis.mimetype === 'text/plain') {
        return content.analysis.body
      } else if (content.analysis.mimetype === 'text/html') {
        return <div
        className={`text-overflow-ellipsis ${styles.emailContent}`}
        dangerouslySetInnerHTML={{ __html: sanitizeHtml(content.analysis.body) }}
      ></div>
      }  else if (content.analysis.mimetype === 'application/x-openline-transcript') {
        try {
          let response = JSON.parse(content.analysis.body);
          return <div>{response.map((transcriptElement: TranscriptElement) => (
             <div><b>{transcriptElement.party}:  </b>{transcriptElement.text}</div>
          ))}</div>;
        } catch (e) {
        }

      }  
    }
    return "Unknown Content: " + JSON.stringify(content);
  }

  const content = decodeContent(message.content);
  if (content.analysis && message.direction) {
    message.direction = 3
  }

  return (
    <>
      {content.dialog && message.direction == 0 && (
        <>
          {(index === 0 ||
            (index > 0 && previousMessage !== message?.direction)) && (
            <div className='mb-1 text-gray-600'>
              {sender.firstName && sender.lastName && (
                <>
                  {sender.firstName} {sender.lastName}
                </>
              )}
              {!sender.firstName && !sender.lastName &&   (
                <>{sender.email}</>
              )}
              {!sender.firstName && !sender.lastName &&  sender.phoneNumber && (
                <>{sender.phoneNumber}</>
              )}
            </div>
          )}
          <div className={styles.messageContainer}>
            <div className={`${styles.message} ${styles.left}`}>
              {displayContent(content)}
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
      {content.dialog && message.direction == 1 && (
        <>
          {(index === 0 ||
            (index > 0 && previousMessage !== message?.direction)) && (
            <div className='w-full flex'>
              <div className='flex-grow-1'></div>
              {<div className="flex-grow-0 mb-1 pr-3">
              {sender.firstName && sender.lastName && (
                <>
                  {sender.firstName} {sender.lastName}
                </>
              )}
              {!sender.firstName && !sender.lastName &&  sender.email && (
                <>{sender.email}</>
              )}
              {!sender.firstName && !sender.lastName &&  sender.phoneNumber && (
                <>{sender.phoneNumber}</>
              )}
              </div>}
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
              {displayContent(content)}
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
      {content.analysis && (
        <>
            <div className='w-full flex'>
              <div className='flex-grow-1'></div>
              {<div className="flex-grow-0 mb-1 pr-3">
              
                <div style={{textAlign: 'center'}}><b>{content.analysis.type}</b></div>
              </div>}
            </div>

   
            <div
              className={`${styles.message} ${styles.center}`}
              style={{ background: '#C5EDCE', borderRadius: '5px' }}
            >
              {displayContent(content)}
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
        </>
      )}
    </>
  );
};

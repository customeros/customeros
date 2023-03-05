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
    parties: Array<VConParty>;
    dialog?: Content
    analysis?: Content
  }

  interface VConParty {
    tel?: string;
    stir?: string;
    mailto?: string;
    name?: string;
  }

  interface TranscriptElement {
    party: VConParty;
    text: string;
  }

  const getUser = (msg: MiniVcon) : VConParty => {
    if (msg.parties) {
      for (const party of msg.parties) {
        if (party.mailto) {
          return party;
        }
      }
    }
    return {mailto: "unknown"};
  }

  const getContact = (msg: MiniVcon) : VConParty  => {
    console.log("getContact"+JSON.stringify(msg)+"\n")
    if (msg.parties) {
      for (const party of msg.parties) {
        if (party.tel) {
          return party;
        }
      }
    }
    return {tel: "unknown"};
  }

  const decodeContent = (content: string)  => {
    let response;
    try {
      response = JSON.parse(content);
    } catch (e) {
      response = {dialog: 
       {
          type: 'MESSAGE',
          mimetype: 'text/plain',
          body: content
       } };
    }
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
          const response = JSON.parse(content.analysis.body);
          return <div>{response.map(
            (transcriptElement: TranscriptElement, index: number) => (
              <div key={index}>
              {transcriptElement.party.tel &&
                (
                  <div
                  className={`${styles.message} ${styles.left}`}>
                  {transcriptElement.text}
                </div>
                )
              }
              {!transcriptElement.party.tel &&
                (
                  (
                    <div
                    className={`${styles.message} ${styles.right}`}
                    style={{ background: '#C5EDCE', borderRadius: '5px' }}
                  >
                    {transcriptElement.text}
                  </div>
                  )
                )
              }
              </div>
          ))}</div>;
        } catch (e) {
          console.log("Got an error: " + e + " when parsing: " + content.analysis.body);
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
              
                <table width="100%"><tr>
                <td><div>{getContact(content).tel}</div></td>
                <td><div style={{textAlign: 'center'}}><b>{content.analysis.type}</b></div></td>
                <td align="right"><div>{getUser(content).mailto}</div></td>
                </tr></table>
              </div>
              }
              
            </div>

   
            <div
              className={`${styles.message} ${styles.center}`}
              style={{ background: '#E5FAE9', borderRadius: '5px' }}
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

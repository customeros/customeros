import React, { useRef, useState } from 'react';
import sanitizeHtml from 'sanitize-html';
import styles from './email-timeline-item.module.scss';
import { Button, IconButton, Reply, ReplyLeft, ReplyMany } from '../../atoms';
import linkifyHtml from 'linkify-html';
import { EmailParticipants } from './email-participants';
import classNames from 'classnames';

interface Props {
  content: string;
  contentType: string;
  sentBy: Array<any>;
  sentTo: Array<any>;
  interactionSession: any;
}

// todo personal section
// email Section

export const EmailTimelineItem: React.FC<Props> = ({
  content,
  contentType,
  sentBy,
  sentTo,
  interactionSession,
  ...rest
}) => {
  let from = '';
  if (
    sentBy &&
    sentBy.length > 0 &&
    sentBy[0].__typename === 'EmailParticipant' &&
    sentBy[0].emailParticipant
  ) {
    from = sentBy[0].emailParticipant.email;
  }

  let to = '';
  if (sentTo && sentTo.length > 0) {
    to = sentTo
      .filter((p: any) => p.type === 'TO')
      .map((p: any) => {
        if (p.__typename === 'EmailParticipant' && p.emailParticipant) {
          return p.emailParticipant.email;
        }
        return '';
      })
      .join('; ');
  }

  let cc = '';
  if (sentTo && sentTo.length > 0) {
    cc = sentTo
      .filter((p: any) => p.type === 'CC')
      .map((p: any) => {
        if (p.__typename === 'EmailParticipant' && p.emailParticipant) {
          return p.emailParticipant.email;
        }
        return '';
      })
      .join('; ');
  }

  let bcc = '';
  if (sentTo && sentTo.length > 0) {
    bcc = sentTo
      .filter((p: any) => p.type === 'BCC')
      .map((p: any) => {
        if (p.__typename === 'EmailParticipant' && p.emailParticipant) {
          return p.emailParticipant.email;
        }
        return '';
      })
      .join('; ');
  }

  const [expanded, toggleExpanded] = useState(false);
  const timelineItemRef = useRef<HTMLDivElement>(null);

  const handleToggleExpanded = () => {
    toggleExpanded(!expanded);
    if (timelineItemRef?.current && expanded) {
      timelineItemRef?.current?.scrollIntoView({ behavior: 'smooth' });
    }
  };

  // const setEditorMode = useSetRecoilState(editorMode);
  // const [emailEditorData, setEmailEditorData] = useRecoilState(editorEmail);
  // const loggedInUserData = useRecoilValue(userData);
  //
  // useEffect(() => {
  //   return () => {
  //     setEmailEditorData({
  //       //@ts-expect-error fixme later
  //       handleSubmit: () => null,
  //       to: [],
  //       subject: '',
  //       respondTo: '',
  //     });
  //     setEditorMode({
  //       mode: EditorMode.Note,
  //       submitButtonLabel: 'Log into timeline',
  //     });
  //   };
  // }, []);

  // const handleSendMessage = (
  //     text: string,
  //     onSuccess: () => void,
  //     destination = [],
  //     replyTo: null | string,
  // ) => {
  //   if (!text) return;
  //   const message: FeedPostRequest = {
  //     channel: 'EMAIL',
  //     username: loggedInUserData.identity,
  //     message: text,
  //     direction: 'OUTBOUND',
  //     destination: destination,
  //   };
  //   if (replyTo) {
  //     message.replyTo = replyTo;
  //   }

  //   axios
  //       .post(`/oasis-api/feed/${feedId}/item`, message)
  //       .then((res) => {
  //         console.log(res);
  //         if (res.data) {
  //           setMessages((messageList: any) => [...messageList, res.data]);
  //           onSuccess();
  //           setEditorMode({
  //             submitButtonLabel: 'Log into timeline',
  //             mode: EditorMode.Note,
  //           });
  //           setEmailEditorData({ ...emailEditorData, to: [], subject: '' });
  //           toast.success('Email sent!');
  //         }
  //       })
  //       .catch(() => {
  //         toast.error('Something went wrong while sending email');
  //       });
  // };
  return (
    <div>
      <div
        className={classNames(styles.emailWrapper, {
          [styles.expanded]: expanded,
        })}
      >
        <div ref={timelineItemRef} className={styles.scrollToView} />
        <article className={`${styles.emailContainer}`}>
          <div>
            <EmailParticipants
              from={from}
              to={to}
              subject={interactionSession.name}
              cc={cc}
              bcc={bcc}
            />
          </div>

          <div
            className={`${styles.emailContentContainer} ${
              !expanded ? styles.eclipse : ''
            }`}
          >
            {contentType === 'text/html' && (
              <div
                className={`text-overflow-ellipsis ${styles.emailContent}`}
                dangerouslySetInnerHTML={{
                  __html: sanitizeHtml(
                    linkifyHtml(content, {
                      defaultProtocol: 'https',
                      rel: 'noopener noreferrer',
                    }),
                  ),
                }}
              ></div>
            )}
            {contentType === 'text/plain' && (
              <div className={`text-overflow-ellipsis ${styles.emailContent}`}>
                {content}
              </div>
            )}

            {!expanded && <div className={styles.eclipse} />}
          </div>
        </article>
      </div>
      <div className={styles.folderTab}>
        <Button
          onClick={() => handleToggleExpanded()}
          mode='link'
          className={styles.toggleExpandButton}
        >
          {expanded ? 'Collapse' : 'Expand'}
        </Button>

        {/*TODO enable after backend refactor*/}
        {/*<Button*/}
        {/*    mode='link'*/}
        {/*    onClick={() => {*/}
        {/*      // TODO add cc and bcc*/}

        {/*      setEmailEditorData({*/}
        {/*        //@ts-expect-error fixme later*/}
        {/*        handleSubmit: handleSendMessage,*/}
        {/*        to: [emailData.from],*/}
        {/*        subject: emailData.subject,*/}
        {/*        respondTo:*/}
        {/*        //@ts-expect-error fixme later*/}
        {/*            msg?.messageId?.conversationEventId || null,*/}
        {/*      });*/}
        {/*      setEditorMode({*/}
        {/*        mode: EditorMode.Email,*/}
        {/*        submitButtonLabel: 'Send',*/}
        {/*      });*/}
        {/*    }}*/}
        {/*>*/}
        {/*  Respond*/}
        {/*</Button>*/}
      </div>
    </div>
  );
};

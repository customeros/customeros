import React, { useState } from 'react';
import styles from './email-participants.module.scss';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { Button } from '@spaces/atoms/button';
import { Avatar } from '@spaces/atoms/avatar';
import { default as Reply } from '@spaces/atoms/icons/Reply';
import { default as ReplyLeft } from '@spaces/atoms/icons/ReplyLeft';
import { default as ReplyMany } from '@spaces/atoms/icons/ReplyMany';
import classNames from 'classnames';
import { SendMailRequest } from '../../conversation-timeline-item/types';
import { useRecoilState, useRecoilValue, useSetRecoilState } from 'recoil';
import {
  editorEmail,
  editorMode,
  EditorMode,
  userData,
} from '../../../../../state';
import axios from 'axios';
import { toast } from 'react-toastify';
import { showLegacyEditor } from '../../../../../state/editor';

interface Props {
  from: string;
  fromName: string;
  to: Array<string>;
  subject: string;
  cc: string;
  bcc: string;
}

// todo personal section
// email Section

export const EmailParticipants: React.FC<Props> = ({
  from,
  to,
  cc,
  bcc,
  subject,
  fromName,
}) => {
  const name = fromName?.split(' ');
  const [showMore, setShowMore] = useState(false);
  const setEditorMode = useSetRecoilState(editorMode);
  const [emailEditorData, setEmailEditorData] = useRecoilState(editorEmail);
  const setShowLegacyEditor = useSetRecoilState(showLegacyEditor);
  const loggedInUserData = useRecoilValue(userData);
  console.log('🏷️ ----- loggedInUserData: '
      , loggedInUserData);
  const SendMail = (
    text: string,
    onSuccess: () => void,
    destination: Array<string> = [],
    replyTo: null | string,
    subject: null | string,
  ) => {
    if (!text) return;
    const request: SendMailRequest = {
      channel: 'EMAIL',
      username: loggedInUserData.identity,
      content: text,
      direction: 'OUTBOUND',
      destination: destination,
    };
    if (replyTo) {
      request.replyTo = replyTo;
    }
    if (subject) {
      request.subject = subject;
    }
    axios
      .post(`/comms-api/mail/send`, request, {
        headers: {
          'X-Openline-Mail-Api-Key': `${process.env.COMMS_MAIL_API_KEY}`,
        },
      })
      .then((res) => {
        if (res.data) {
          onSuccess();
          setShowLegacyEditor(false);
          setEditorMode({
            mode: EditorMode.Note,
          });
          setEmailEditorData({ ...emailEditorData, to: [], subject: '' });
        }
      })
      .catch(() => {
        toast.error('Something went wrong while sending request');
      });
  };
  return (
    <div className={styles.wrapper}>
      <section className={styles.emailDataContainer}>
        <div className={styles.avatar}>
          <Avatar name={name[0]} surname={name?.[1]} size={30} />
        </div>
        <div className='flex w-full flex-column'>
          <div className={classNames(styles.emailDataRow, styles.participants)}>
            <div className={styles.emailParticipantWrapper}>
              <div
                className={classNames(styles.label, styles.labelWithSpacing)}
              >
                From:
              </div>
              <div className={styles.data}>{from}</div>
              <div
                className={classNames(
                  styles.emailParticipantWrapper,
                  styles.nowrap,
                )}
              >
                <div
                  className={classNames(styles.label, {
                    [styles.labelWithSpacing]: to?.length > 1 || to.length > 10,
                  })}
                >
                  To:
                </div>
                <div className={styles.data}>{to && to.join(',')}</div>
              </div>
            </div>
          </div>
          {showMore && (
            <div>
              {!!cc?.length && (
                <div className={styles.emailDataRow}>
                  <div
                    className={classNames(
                      styles.label,
                      styles.labelWithSpacing,
                    )}
                  >
                    Cc:
                  </div>
                  <div className={styles.data}>{cc.split(';').join(',')}</div>
                </div>
              )}
              {!!bcc?.length && (
                <div className={styles.emailDataRow}>
                  <div
                    className={classNames(
                      styles.label,
                      styles.labelWithSpacing,
                    )}
                  >
                    Bcc:
                  </div>
                  <div className={styles.data}>{bcc.split(';').join(',')}</div>
                </div>
              )}
            </div>
          )}
          <div>
            <div className={styles.emailDataRow}>
              <div
                className={classNames(styles.label, styles.labelWithSpacing)}
              >
                Subject:
              </div>
              <div className={styles.data}>{subject}</div>
            </div>
          </div>
        </div>
      </section>
      <div className={styles.buttonsContainer}>
        <Button
          mode='link'
          className={styles.showMoreBtn}
          onClick={() => setShowMore(!showMore)}
        >
          Cc / Bcc
        </Button>

        <IconButton
          className={styles.emailActionButton}
          isSquare
          label='Reply'
          size='xxxs'
          onClick={() => {
            setEmailEditorData({
              //@ts-expect-error fixme later
              handleSubmit: (data) => {
                SendMail(data, () => null, [from], null, subject);
              },
              to: [from],
              subject: subject,
              respondTo: from,
            });
            setShowLegacyEditor(true);
            setEditorMode({
              mode: EditorMode.Email,
            });
          }}
          icon={<ReplyLeft />}
        />
        <IconButton
          label='Reply many'
          className={styles.emailActionButton}
          isSquare
          size='xxxs'
          onClick={() => null}
          disabled
          icon={<ReplyMany />}
        />
        <IconButton
          label=''
          className={styles.emailActionButton}
          isSquare
          size='xxxs'
          onClick={() => null}
          disabled
          icon={<Reply />}
        />
      </div>
    </div>
  );
};

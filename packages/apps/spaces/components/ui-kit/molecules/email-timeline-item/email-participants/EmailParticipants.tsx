import React, { useState } from 'react';
import styles from './email-participants.module.scss';
import { useContactNameFromEmail } from '../../../../../hooks/useContact';
import {
  Avatar,
  Button,
  IconButton,
  Reply,
  ReplyLeft,
  ReplyMany,
} from '../../../atoms';
import { getContactDisplayName } from '../../../../../utils';
import classNames from 'classnames';

interface Props {
  from: string;
  to: string;
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
}) => {
  const { loading, error, data } = useContactNameFromEmail({ email: from });
  const name = getContactDisplayName(data)?.split(' ');
  const [showMore, setShowMore] = useState(false);

  return (
    <div className={styles.wrapper}>
      <section className={styles.emailDataContainer}>
        {!loading && !error && (
          <div className={styles.avatar}>
            <Avatar
              name={name?.[0] || ''}
              surname={name.length === 2 ? name[1] : name[2]}
              size={30}
            />
          </div>
        )}
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
                    [styles.labelWithSpacing]: to?.split(';').length > 1,
                  })}
                >
                  To:
                </div>
                <div className={styles.data}>{to.split(';').join(',')}</div>
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
          size='xxxs'
          onClick={() => null}
          disabled
          icon={<ReplyLeft />}
        />
        <IconButton
          className={styles.emailActionButton}
          isSquare
          size='xxxs'
          onClick={() => null}
          disabled
          icon={<ReplyMany />}
        />
        <IconButton
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

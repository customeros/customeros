import React from 'react';
import { Avatar, User } from '../../atoms';
import styles from './conversation-timeline-item.module.scss';
import {
  useContactNameFromEmail,
  useContactNameFromPhoneNumber,
} from '../../../../hooks/useContact';
import { getContactDisplayName } from '../../../../utils';
import classNames from 'classnames';

interface Props {
  direction: number;
  sender: any;
}
export const CallParties: React.FC<Props> = ({ direction, sender }) => {
  const { data } = useContactNameFromEmail({
    email: sender?.senderUsername.identifier || '',
  });
  const { data: data1 } = useContactNameFromPhoneNumber({
    phoneNumber: sender?.senderUsername.identifier || '',
  });

  const getName = () => {
    return (data || data1) && getContactDisplayName(data || data1) !== 'Unnamed'
      ? getContactDisplayName(data || data1)
      : undefined;
  };

  const initials = getContactDisplayName(data || data1).split(' ');

  return (
    <div
      className={classNames(styles.avatarSection, {
        [styles.directionRight]: direction === 1,
      })}
    >
      <Avatar
        name={initials[0]}
        surname={initials.length === 2 ? initials[1] : initials[2]}
        size={30}
        image={
          initials.length === 1 ? (
            <User style={{ transform: 'scale(0.8)' }} />
          ) : undefined
        }
      />
      <div className={styles.contactName}>
        {getName() || sender?.senderUsername.identifier}
      </div>
    </div>
  );
};

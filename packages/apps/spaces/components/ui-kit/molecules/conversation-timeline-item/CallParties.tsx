import React, { useEffect, useState } from 'react';
import { Avatar } from '@spaces/atoms/avatar';
import User from '@spaces/atoms/icons/User';
import styles from './conversation-timeline-item.module.scss';
import {
  useContactNameFromEmail,
  useContactNameFromId,
  useContactNameFromPhoneNumber,
} from '@spaces/hooks/useContact';
import { getContactDisplayName } from '../../../../utils';
import classNames from 'classnames';

interface Props {
  direction: number;
  sender: any;
  mode: 'PHONE_CALL' | 'CHAT' | 'LIVE';
}
export const CallParties: React.FC<Props> = ({ direction, sender, mode }) => {
  const { data: dataFromEmail } = useContactNameFromEmail({
    email: sender?.senderUsername?.identifier || '',
  });
  const { data: dataFromPhoneNumber } = useContactNameFromPhoneNumber({
    phoneNumber: mode !== 'CHAT' ? sender?.senderUsername.identifier : '',
  });
  const { data: dataFromId } = useContactNameFromId({
    id: mode !== 'CHAT' ? sender?.senderId : sender?.senderUsername?.identifier,
  });
  const [initials, setInitials] = useState<Array<string>>([]);
  const [name, setName] = useState('');

  useEffect(() => {
    const data = dataFromId || dataFromEmail || dataFromPhoneNumber;

    if (
      !initials.length &&
      (sender?.senderUsername?.identifier || sender?.senderId) &&
      data
    ) {
      const name = getContactDisplayName(data);
      const initials = (name !== 'Unnamed' ? name : '').split(' ');
      if (initials.length) {
        setInitials(initials);
      }
      setName(name || sender?.senderUsername?.identifier);
    }
  }, [
    dataFromId,
    dataFromEmail,
    dataFromPhoneNumber,
    sender?.senderId,
    sender?.senderUsername?.identifier,
  ]);

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
          initials.length < 2 ? (
            <User height={20} />
          ) : undefined
        }
      />
      <div className={styles.contactName}>
        {name || sender?.senderUsername?.identifier || 'Unnamed'}
      </div>
    </div>
  );
};

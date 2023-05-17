import React, { useEffect, useState } from 'react';
import { Avatar } from '@spaces/atoms/avatar';
import User from '@spaces/atoms/icons/User';
import styles from './conversation-timeline-item.module.scss';
import classNames from 'classnames';
import { useContactNameLazy } from '@spaces/hooks/useContact/useContactNameLazy';

interface Props {
  direction: number;
  sender: any;
  mode: 'PHONE_CALL' | 'CHAT' | 'LIVE';
}
export const CallParties: React.FC<Props> = ({ direction, sender, mode }) => {
  const { contactName } = useContactNameLazy({
    email: sender?.senderUsername?.identifier,
    phoneNumber: mode !== 'CHAT' ? sender?.senderUsername.identifier : '',
    id: mode !== 'CHAT' ? sender?.senderId : sender?.senderUsername?.identifier,
  });

  const [initials, setInitials] = useState<Array<string>>([]);

  useEffect(() => {
    if (
      !initials.length &&
      (sender?.senderUsername?.identifier || sender?.senderId) &&
      contactName
    ) {
      const initials = (contactName !== 'Unnamed' ? contactName : '').split(
        ' ',
      );
      if (initials.length) {
        setInitials(initials);
      }
    }
  }, [contactName, sender?.senderId, sender?.senderUsername?.identifier]);

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
        image={initials.length < 2 ? <User height={20} /> : undefined}
      />
      <div className={styles.contactName}>
        {contactName || sender?.senderUsername?.identifier || 'Unnamed'}
      </div>
    </div>
  );
};

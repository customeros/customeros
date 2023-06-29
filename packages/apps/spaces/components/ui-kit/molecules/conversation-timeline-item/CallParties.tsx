import React, { useEffect, useState } from 'react';
import { Avatar } from '@spaces/atoms/avatar';
import User from '@spaces/atoms/icons/User';
import styles from './conversation-timeline-item.module.scss';
import classNames from 'classnames';

interface Props {
  direction: number;
  name: string;
}
export const CallParties: React.FC<Props> = ({ direction, name }) => {
  const [initials, setInitials] = useState<Array<string>>([]);

  useEffect(() => {
    if (!initials.length && name) {
      const initials = (name !== 'Unnamed' ? name : '').split(' ');
      if (initials.length) {
        setInitials(initials);
      }
    }
  }, [initials.length, name]);

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
      <div className={styles.contactName}>{name}</div>
    </div>
  );
};

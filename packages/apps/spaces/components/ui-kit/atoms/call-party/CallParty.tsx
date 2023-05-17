import React, { useEffect, useState } from 'react';
import { Avatar, User } from '../index';
import styles from './call-party.module.scss';

import classNames from 'classnames';

interface Props {
  direction: 'right' | 'left';
  name: string;
}
export const CallParty: React.FC<Props> = ({ direction, name }) => {
  const [initials, setInitials] = useState<Array<string>>([]);

  useEffect(() => {
    if (name) {
      const initialsFromName = (name !== 'Unnamed' ? name : '').split(' ');
      if (initialsFromName.length) {
        setInitials(initialsFromName);
      }
    }
  }, [name]);

  return (
    <div
      className={classNames(styles.avatarSection, {
        [styles.directionRight]: direction === 'right',
      })}
    >
      <Avatar
        name={initials[0]}
        surname={initials.length === 2 ? initials[1] : initials[2]}
        size={30}
        image={initials.length < 2 ? <User height={16} /> : undefined}
      />
      <div className={classNames(styles.contactName)}>{name || 'Unnamed'}</div>
    </div>
  );
};

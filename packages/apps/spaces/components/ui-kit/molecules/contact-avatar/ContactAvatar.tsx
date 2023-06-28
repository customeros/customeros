import React, { memo, useEffect, useState } from 'react';
import { default as User } from '../../atoms/icons/User';
import { Avatar } from '../../atoms/avatar';

interface Props {
  name: string;
  size?: number;
  showName?: boolean;
  onlyName?: boolean;
}

export const ContactAvatar: React.FC<Props> = memo(
  function ContactAvatarComponent({
    showName = false,
    onlyName = false,
    size = 24,
    name,
  }) {
    const [contactName, setContactName] = useState(['', '']);

    useEffect(() => {
      setContactName(name.split(' '));
    }, [name]);

    return (
      <>
        {!onlyName && (
          <Avatar
            name={contactName?.[0] || ''}
            surname={contactName.length === 2 ? contactName[1] : contactName[2]}
            size={size}
            image={contactName.length === 1 && <User />}
          />
        )}

        {(showName || onlyName) && <div>{contactName}</div>}
      </>
    );
  },
  (prevProps, nextProps) =>
    prevProps.name === nextProps.name &&
    nextProps.size === prevProps.size,
);

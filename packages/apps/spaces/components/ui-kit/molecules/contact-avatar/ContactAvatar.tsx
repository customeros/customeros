import React, { memo } from 'react';
import { default as User } from '../../atoms/icons/User';
import { Avatar } from '../../atoms/avatar';
import { useContactNameFromId } from '@spaces/hooks/useContact';
import { getContactDisplayName } from '../../../../utils';
import { Skeleton } from '../../atoms/skeleton';

interface Props {
  contactId: string;
  size?: number;
  showName?: boolean;
  onlyName?: boolean;
}

export const ContactAvatar: React.FC<Props> = memo(
  function ContactAvatarComponent({
    contactId,
    showName = false,
    onlyName = false,
    size = 30,
  }) {
    const { loading, error, data } = useContactNameFromId({ id: contactId });
    if (loading || error) {
      if (showName) {
        return <Skeleton />;
      }

      return <div />;
    }
    const name = getContactDisplayName(data).split(' ');
    console.log('🏷️ ----- name: ', name);
    return (
      <>
        {!onlyName && (
          <Avatar
            name={name?.[0] || ''}
            surname={name.length === 2 ? name[1] : name[2]}
            size={size}
            image={name.length === 1 && <User />}
          />
        )}

        {(showName || onlyName) && <div>{name}</div>}
      </>
    );
  },
  (prevProps, nextProps) =>
    prevProps.contactId === nextProps.contactId &&
    nextProps.size === prevProps.size,
);

import React, { memo, useEffect, useState } from 'react';
import { default as User } from '../../atoms/icons/User';
import { Avatar } from '../../atoms/avatar';
import { useContactNameFromId } from '@spaces/hooks/useContact';
import { getContactDisplayName } from '../../../../utils';
import { Skeleton } from '../../atoms/skeleton';

interface Props {
  contactId: string;
  name?: string;
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
    name = '',
  }) {
    const { onGetContactNameById, loading, error } = useContactNameFromId();
    const [contactName, setContactName] = useState(['', '']);

    const handleGetContactNameById = async () => {
      const result = await onGetContactNameById({
        variables: { id: contactId },
      });
      if (result.name || result.firstName || result.lastName) {
        setContactName(getContactDisplayName(result).split(' '));
      }
    };

    useEffect(() => {
      if (!name) {
        handleGetContactNameById();
      }
      if (name) {
        setContactName(name.split(' '));
      }
    }, [name]);

    if (loading || error) {
      if (showName) {
        return <Skeleton />;
      }

      return <div />;
    }

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
    prevProps.contactId === nextProps.contactId &&
    prevProps.name === nextProps.name &&
    nextProps.size === prevProps.size,
);

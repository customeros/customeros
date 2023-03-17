import React from 'react';
import { Avatar } from '../../atoms';
import { useContactName } from '../../../../hooks/useContact';
import { getContactDisplayName } from '../../../../utils';

interface Props {
  contactId: string;
  size?: number;
}

export const ContactAvatar: React.FC<Props> = ({ contactId, size = 30 }) => {
  const { loading, error, data } = useContactName({ id: contactId });
  if (loading || error) {
    return <div />;
  }
  const name = getContactDisplayName(data).split(' ');
  console.log('🏷️ ----- name: ', name);
  return (
    <Avatar
      name={name?.[0] || ''}
      surname={name.length === 2 ? name[1] : name[2]}
      size={size}
    />
  );
};

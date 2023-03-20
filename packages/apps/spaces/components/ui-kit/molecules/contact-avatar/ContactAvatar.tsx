import React from 'react';
import { Avatar } from '../../atoms';
import { useContactNameFromId } from '../../../../hooks/useContact';
import { getContactDisplayName } from '../../../../utils';

interface Props {
  contactId: string;
  size?: number;
}

export const ContactAvatar: React.FC<Props> = ({ contactId, size = 30 }) => {
  const { loading, error, data } = useContactNameFromId({ id: contactId });
  if (loading || error) {
    return <div />;
  }
  const name = getContactDisplayName(data).split(' ');
  return (
    <Avatar
      name={name?.[0] || ''}
      surname={name.length === 2 ? name[1] : name[2]}
      size={size}
    />
  );
};

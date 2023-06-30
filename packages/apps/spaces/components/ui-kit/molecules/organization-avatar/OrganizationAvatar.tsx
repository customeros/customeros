import React from 'react';
import { Avatar } from '@spaces/atoms/avatar';
import { Building } from '@spaces/atoms/icons';

interface Props {
  name?: string;
  size?: number;
}

export const OrganizationAvatar: React.FC<Props> = ({
  size = 24,
  name = '',
}) => {
  return (
    <Avatar
      name={name}
      surname={''}
      size={size}
      image={!name && <Building />}
      isSquare
    />
  );
};

import React from 'react';
import { Avatar } from '../../atoms';
import { useOrganizationName } from '../../../../hooks/useOrganization/useOrganizationName';

interface Props {
  organizationId: string;
  size?: number;
}

export const OrganizationAvatar: React.FC<Props> = ({
  organizationId,
  size = 35,
}) => {
  const { loading, error, data } = useOrganizationName({ id: organizationId });
  if (loading || error) {
    return <div />;
  }
  const name = (data?.name ?? '').split(' ');
  return (
    <Avatar
      name={data?.name || ''}
      surname={name?.length > 1 ? name?.[1] : ''}
      size={size}
      isSquare
    />
  );
};

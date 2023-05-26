import React, { useEffect, useState } from 'react';
import { Avatar } from '@spaces/atoms/avatar';
import { useOrganizationName } from '@spaces/hooks/useOrganization/useOrganizationName';
import { Building } from '@spaces/atoms/icons';

interface Props {
  organizationId: string;
  name?: string;
  size?: number;
}

export const OrganizationAvatar: React.FC<Props> = ({
  organizationId,
  size = 30,
  name = '',
}) => {
  const { loading, error, onGetOrganizationName } = useOrganizationName();
  const [organizationName, setOrganizationName] = useState(name);

  const handleGetOrganizationNameById = async () => {
    const result = await onGetOrganizationName({
      variables: { id: organizationId },
    });
    if (result.name) {
      setOrganizationName(result.name);
    }
  };

  useEffect(() => {
    if (!name) {
      handleGetOrganizationNameById();
    }
  }, [name]);

  if (loading || error) {
    return <div />;
  }
  return (
    <Avatar
      name={organizationName}
      surname={''}
      size={size}
      image={!organizationName && <Building />}
      isSquare
    />
  );
};

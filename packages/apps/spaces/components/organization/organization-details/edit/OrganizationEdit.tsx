import React from 'react';
import { OrganizationForm } from './OrganizationForm';
import { useRouter } from 'next/router';
import { useUpdateOrganization } from '../../../../hooks/useOrganization/useUpdateOrganization';

export const OrganizationEdit = ({
  data,
  onSetMode,
}: {
  data: any;
  onSetMode: any;
}) => {
  const { onUpdateOrganization } = useUpdateOrganization({
    organizationId: data.id,
  });
  const handleCreateOrganization = (values: any) => {
    onUpdateOrganization(values).then((value: any) => {
      if (value?.id) {
        onSetMode('PREVIEW');
      }
    });
  };
  return <OrganizationForm data={data} onSubmit={handleCreateOrganization} />;
};

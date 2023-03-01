import React from 'react';
import { useCreateOrganization } from '../../../../hooks/useOrganization';
import { OrganizationForm } from './OrganizationForm';
import { useRouter } from 'next/router';

export const OrganizationCreate: React.FC = () => {
  const router = useRouter();

  const { onCreateOrganization } = useCreateOrganization();
  const handleCreateOrganization = (values: any) => {
    onCreateOrganization(values).then((value) => {
      if (value?.id) {
        router.push(`/organization/${value.id}`);
      }
    });
  };
  return <OrganizationForm onSubmit={handleCreateOrganization} mode='CREATE' />;
};

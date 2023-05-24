import React from 'react';
import {
  useOrganizationLocations,
  useCreateOrganizationLocation,
} from '@spaces/hooks/useOrganizationLocation';
import { LocationList } from '../../shared/location';

interface OrganizationLocationsProps {
  id: string;
}

export const OrganizationLocations: React.FC<OrganizationLocationsProps> = ({
  id,
}) => {
  const { data, loading, error } = useOrganizationLocations({ id });
  const { onCreateOrganizationLocation } = useCreateOrganizationLocation({
    organizationId: id,
  });

  if (loading) return null;
  if (error) {
    return (
      <div>Sorry looks like there was an error during loading locations</div>
    );
  }
  return (
    <LocationList
      locations={data?.locations || []}
      onCreateLocation={onCreateOrganizationLocation}
    />
  );
};

import React from 'react';
import { useCreateOrganizationLocation } from '@spaces/hooks/useOrganizationLocation';
import { LocationList } from '../../shared/location';
import { useRecoilValue } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { LocationListSkeleton } from '../../shared/location/skeletons/LocationListSkeleton';

interface OrganizationLocationsProps {
  id: string;
  loading: boolean;
  locations: Array<any> | undefined | null;
}

export const OrganizationLocations: React.FC<OrganizationLocationsProps> = ({
  id,
  locations,
  loading,
}) => {
  const { isEditMode } = useRecoilValue(organizationDetailsEdit);
  const { onCreateOrganizationLocation } = useCreateOrganizationLocation({
    organizationId: id,
  });

  if (loading) {
    return <LocationListSkeleton />;
  }
  return (
    <LocationList
      isEditMode={isEditMode}
      locations={locations || []}
      onCreateLocation={onCreateOrganizationLocation}
    />
  );
};

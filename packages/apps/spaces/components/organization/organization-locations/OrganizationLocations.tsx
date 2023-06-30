import React from 'react';
import { useCreateOrganizationLocation } from '@spaces/hooks/useOrganizationLocation';
import { LocationList } from '../../shared/location';
import { useRecoilValue } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { LocationListSkeleton } from '../../shared/location/skeletons/LocationListSkeleton';
import { useRemoveOrganizationLocation } from '@spaces/hooks/useOrganizationLocation/useRemoveOrganizationLocation';
import { Location as COSLocation } from '@spaces/graphql';

type TLocation = Omit<
  COSLocation,
  'appSource' | 'source' | 'sourceOfTruth' | 'createdAt' | 'updatedAt'
>;
interface OrganizationLocationsProps {
  id: string;
  loading: boolean;
  locations: Array<TLocation> | undefined | null;
}

export const OrganizationLocations: React.FC<OrganizationLocationsProps> = ({
  id,
  locations,
  loading,
}) => {
  const { isEditMode } = useRecoilValue(organizationDetailsEdit);
  const { onRemoveOrganizationLocation } = useRemoveOrganizationLocation({
    organizationId: id,
  });
  const { onCreateOrganizationLocation } = useCreateOrganizationLocation({
    organizationId: id,
  });

  if (loading) {
    return <LocationListSkeleton />;
  }
  return (
    <LocationList
      isEditMode={isEditMode}
      onRemoveLocation={onRemoveOrganizationLocation}
      locations={locations || []}
      onCreateLocation={onCreateOrganizationLocation}
    />
  );
};

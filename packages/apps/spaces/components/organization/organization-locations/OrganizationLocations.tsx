import React from 'react';
import {
  useOrganizationLocations,
  useCreateOrganizationLocation,
} from '@spaces/hooks/useOrganizationLocation';
import { LocationList } from '../../shared/location';
import { useRecoilValue } from 'recoil';
import { organizationDetailsEdit } from '../../../state';

interface OrganizationLocationsProps {
  id: string;
}

export const OrganizationLocations: React.FC<OrganizationLocationsProps> = ({
  id,
}) => {
  const { data, error } = useOrganizationLocations({ id });
  const { isEditMode } = useRecoilValue(organizationDetailsEdit);
  const { onCreateOrganizationLocation } = useCreateOrganizationLocation({
    organizationId: id,
  });

  if (error) {
    return (
      <div>Sorry looks like there was an error during loading locations</div>
    );
  }
  return (
    <LocationList
      isEditMode={isEditMode}
      locations={data?.locations || []}
      onCreateLocation={onCreateOrganizationLocation}
    />
  );
};

import React from 'react';
import {
  useCreateContactLocation,
} from '@spaces/hooks/useContactLocation';
import { LocationList } from '../../shared/location';
import { useRecoilValue } from 'recoil';
import { contactDetailsEdit } from '../../../state';
import { LocationListSkeleton } from '../../shared/location/skeletons/LocationListSkeleton';

interface ContactLocationsProps {
  id: string;
  data: any;
  loading: boolean;
}

export const ContactLocations: React.FC<ContactLocationsProps> = ({
  id,
  data,
  loading,
}) => {
  const { isEditMode } = useRecoilValue(contactDetailsEdit);
  const { onCreateContactLocation } = useCreateContactLocation({
    contactId: id,
  });

  if (loading) {
    return <LocationListSkeleton />;
  }

  return (
    <LocationList
      isEditMode={isEditMode}
      locations={data?.locations || []}
      onCreateLocation={onCreateContactLocation}
    />
  );
};

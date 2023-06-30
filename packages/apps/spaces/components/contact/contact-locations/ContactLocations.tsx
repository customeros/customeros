import React from 'react';
import { useCreateContactLocation } from '@spaces/hooks/useContactLocation';
import { LocationList } from '../../shared/location';
import { useRecoilValue } from 'recoil';
import { contactDetailsEdit } from '../../../state';
import { LocationListSkeleton } from '../../shared/location/skeletons/LocationListSkeleton';
import { useRemoveContactLocation } from '@spaces/hooks/useContactLocation/useRemoveContactLocation';
import { ContactResponse } from '@spaces/hooks/useContact/useGetContact';

interface ContactLocationsProps {
  id: string;
  data: ContactResponse;
  loading: boolean;
}

export const ContactLocations: React.FC<ContactLocationsProps> = ({
  id,
  data,
  loading,
}) => {
  const { onRemoveContactLocation } = useRemoveContactLocation({
    contactId: id,
  });
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
      onRemoveLocation={onRemoveContactLocation}
      locations={data?.locations || []}
      onCreateLocation={onCreateContactLocation}
    />
  );
};

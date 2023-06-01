import React from 'react';
import {
  useContactLocations,
  useCreateContactLocation,
} from '@spaces/hooks/useContactLocation';
import { LocationList } from '../../shared/location';
import { useRecoilValue } from 'recoil';
import { contactDetailsEdit } from '../../../state';

interface ContactLocationsProps {
  id: string;
}

export const ContactLocations: React.FC<ContactLocationsProps> = ({ id }) => {
  const { data, loading, error } = useContactLocations({ id });
  const { isEditMode } = useRecoilValue(contactDetailsEdit);
  const { onCreateContactLocation } = useCreateContactLocation({
    contactId: id,
  });

  if (loading) return null;
  if (error) {
    return (
      <div>Sorry looks like there was an error during loading locations</div>
    );
  }
  return (
    <LocationList
      isEditMode={isEditMode}
      locations={data?.locations || []}
      onCreateLocation={onCreateContactLocation}
    />
  );
};

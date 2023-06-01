import React, { useState } from 'react';
import PinAltLight from '@spaces/atoms/icons/PinAltLight';
import useGoogle from 'react-google-autocomplete/lib/usePlacesAutocompleteService';
import { DebouncedAutocomplete } from '@spaces/atoms/autocomplete';
import { useUpdateLocation } from '@spaces/hooks/useUpdateLocation';

export const Location: React.FC<{
  locationId: string;
  rawAddress: string;
  isEditMode: boolean;
}> = ({ locationId, rawAddress, isEditMode }) => {
  const { placePredictions, getPlacePredictions } = useGoogle({
    apiKey: process.env.GOOGLE_MAPS_API_KEY,
  });
  const [value, setValue] = useState(rawAddress);
  const { onUpdateLocation } = useUpdateLocation();
  return (
    <>
      <PinAltLight style={{ marginRight: 8 }} />
      <DebouncedAutocomplete
        initialValue={value}
        onSearch={(filter: string) => {
          getPlacePredictions({ input: filter });
        }}
        onChange={(d) => {
          setValue(d?.value);
          onUpdateLocation({ id: locationId, rawAddress: d?.value });
        }}
        editable={isEditMode}
        disabled={false}
        newItemLabel={''}
        mode='invisible'
        suggestions={(placePredictions || []).map((suggestion) => ({
          label: suggestion.description,
          value: suggestion.description,
        }))}
      />
    </>
  );
};

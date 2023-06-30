import React from 'react';
import PinAltLight from '@spaces/atoms/icons/PinAltLight';
import useGoogle from 'react-google-autocomplete/lib/usePlacesAutocompleteService';
import { useUpdateLocation } from '@spaces/hooks/useUpdateLocation';
import {
  Select,
  SelectMenu,
  SelectInput,
  SelectWrapper,
} from '@spaces/ui/form/select';
import { DeleteIconButton } from '@spaces/ui-kit/atoms';

export const Location: React.FC<{
  locationId: string;
  locationString: string;
  isEditMode: boolean;
  onRemoveLocation: (locationId: string) => void;
}> = ({ locationId, locationString, isEditMode, onRemoveLocation }) => {
  const { placePredictions, getPlacePredictions } = useGoogle({
    apiKey: process.env.GOOGLE_MAPS_API_KEY,
  });
  const { onUpdateLocation, saving } = useUpdateLocation();

  const getFormattedPredictions = (placePredictions || []).map(
    (suggestion) => ({
      label: suggestion.description,
      value: suggestion.description,
    }),
  );

  const existingOptions = locationString
    ? [{ label: locationString, value: locationId }]
    : [];

  return (
    <>
      <Select<string>
        onSelect={(val) =>
          onUpdateLocation({ id: locationId, rawAddress: val })
        }
        onChange={(val) => getPlacePredictions({ input: val })}
        value={locationId}
        options={[...getFormattedPredictions, ...existingOptions]}
      >
        <SelectWrapper>
          {isEditMode ? (
            <DeleteIconButton
              onDelete={() => onRemoveLocation(locationId)}
              style={{ marginRight: 8 }}
            />
          ) : (
            <PinAltLight style={{ marginRight: 8 }} />
          )}
          <SelectInput
            saving={saving}
            placeholder='Location'
            readOnly={!isEditMode}
          />
          {isEditMode && <SelectMenu />}
        </SelectWrapper>
      </Select>
    </>
  );
};

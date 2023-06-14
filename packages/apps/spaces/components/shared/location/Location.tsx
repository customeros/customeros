import React from 'react';
import PinAltLight from '@spaces/atoms/icons/PinAltLight';
import useGoogle from 'react-google-autocomplete/lib/usePlacesAutocompleteService';
import { useUpdateLocation } from '@spaces/hooks/useUpdateLocation';
import { Select, useSelect } from '@spaces/atoms/select';
import { SelectWrapper } from '@spaces/atoms/select/SelectWrapper';
import { SelectInput } from '@spaces/atoms/select/SelectInput';
import styles from '@spaces/atoms/select/select.module.scss';
import classNames from 'classnames';
interface SelectMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
}
const SelectMenu = ({
  noOfVisibleItems = 7,
  itemSize = 38,
}: SelectMenuProps) => {
  const { state, getMenuProps, getMenuItemProps } = useSelect();
  const maxMenuHeight = itemSize * noOfVisibleItems;
  return (
    <ul
      className={styles.dropdownMenu}
      {...getMenuProps({ maxHeight: maxMenuHeight })}
    >
      {state.items.length ? (
        state.items.map(({ value, label }, index) => (
          <li
            key={value}
            className={classNames(styles.dropdownMenuItem, {
              [styles.isFocused]: state.currentIndex === index,
              [styles.isSelected]: state.selection === value,
            })}
            {...getMenuItemProps({ value, index })}
          >
            {label}
          </li>
        ))
      ) : (
        <li />
      )}
    </ul>
  );
};

export const Location: React.FC<{
  locationId: string;
  locationString: string;
  isEditMode: boolean;
}> = ({ locationId, locationString, isEditMode }) => {
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
          <PinAltLight style={{ marginRight: 8 }} />
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

import React, { useState, useRef } from 'react';
import {
  AutoComplete as PrimereactAutocomplete,
  AutoCompleteChangeParams,
} from 'primereact/autocomplete';
import styles from './autocomplete.module.scss';
import classNames from 'classnames';

interface SuggestionItem {
  label: string;
  value: string;
}

interface CustomAutoCompleteProps {
  value: string;
  createItemType?: string;
  suggestions: SuggestionItem[];
  onChange: (value: SuggestionItem) => void;
  onAddNew: (item: SuggestionItem) => Promise<any>;
  newItemLabel: string;
  editable?: boolean;
  disabled?: boolean;
  placeholder?: string;
}

export const Autocomplete = ({
  value,
  suggestions = [],
  onChange,
  onAddNew,
  createItemType = '',
  newItemLabel = '',
  editable,
  disabled,
  placeholder = '',
}: CustomAutoCompleteProps) => {
  const [inputValue, setInputValue] = useState<string>(value);
  const [filteredSuggestions, setFilteredSuggestions] =
    useState<SuggestionItem[]>(suggestions);
  const inputRef = useRef<HTMLInputElement>(null);
  const handleInputChange = (event: AutoCompleteChangeParams) => {
    const newInputValue = event.value;
    setInputValue(newInputValue);

    // // Filter the suggestions based on the entered value
    // const filteredItems = suggestions.filter((item) =>
    //   item.label.toLowerCase().includes(newInputValue.toLowerCase()),
    // );
    // setFilteredSuggestions(filteredItems);
  };

  const handleSelectItem = (event: { value: SuggestionItem }) => {
    const selectedValue = event.value;
    setInputValue(selectedValue.label);
    onChange(selectedValue);
  };

  const handleAddNew = () => {
    const newItem = { label: inputValue, value: inputValue };
    onAddNew(newItem);
    setInputValue('');
    inputRef?.current?.focus();
  };

  const search = (event: any) => {
    const query = event.query;
    const filteredItems = (suggestions || []).filter(
      (item) => item.label.toLowerCase().indexOf(query.toLowerCase()) !== -1,
    );

    setFilteredSuggestions(filteredItems || []);
  };

  const handleCreateItem = async () => {
    try {
      const newItem = await onAddNew({ value: inputValue, label: inputValue });

      if (newItem) {
        handleSelectItem({ label: newItem[newItemLabel], value: newItem.id });
      }
    } catch (e) {
      console.log('🏷️ ----- : ERRROR', e);
    }
  };

  return (
    <div className={styles.autocompleteContainer}>
      <div>
        <PrimereactAutocomplete
          inputClassName={classNames(styles.autocompleteInput, {
            [styles.notEditable]: !editable,
            [styles.disabled]: disabled,
          })}
          disabled={!editable || disabled}
          value={inputValue}
          delay={300}
          placeholder={placeholder}
          suggestions={filteredSuggestions}
          onChange={handleInputChange}
          itemTemplate={(data) => (
            <span onClick={() => handleSelectItem(data)}>{data.label}</span>
          )}
          completeMethod={search}
          onSelect={handleSelectItem}
          onKeyUp={(event) => {
            if (
              event.key === 'Enter' &&
              !filteredSuggestions.find((item) => item.label === inputValue)
            ) {
              handleAddNew();
            }
          }}
          inputRef={inputRef}
        />
      </div>

      {inputValue && inputValue !== value && !filteredSuggestions.length && (
        <div className={styles.createItemButton}>
          <button onClick={handleCreateItem}>
            Create {createItemType} &apos;{inputValue}&apos;
          </button>
        </div>
      )}
    </div>
  );
};

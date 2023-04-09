import React, { useState, useRef, useLayoutEffect, useEffect } from 'react';
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
  mode?: 'default' | 'fit-content';
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
  mode = 'default',
}: CustomAutoCompleteProps) => {
  const [inputValue, setInputValue] = useState<string>(value);
  const [width, setWidth] = useState<number>();
  const [showCreateButton, setShowCreateButton] = useState<boolean>(false);
  const [filteredSuggestions, setFilteredSuggestions] =
    useState<SuggestionItem[]>(suggestions);
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLSpanElement>(null);
  const handleInputChange = (event: AutoCompleteChangeParams) => {
    const newInputValue = event.value;
    setInputValue(newInputValue);
  };

  useLayoutEffect(() => {
    if (mode === 'fit-content') {
      setWidth((measureRef?.current?.scrollWidth || 0) + 2);
    }
  }, [inputValue]);
  useEffect(() => {
    if (inputValue && inputValue !== value && !filteredSuggestions.length) {
      setShowCreateButton(true);
    }
  }, [filteredSuggestions.length, inputValue, value]);

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
        handleSelectItem({
          value: {
            label: newItem[newItemLabel],
            value: newItem.id,
          },
        });
        setShowCreateButton(false);
      }
    } catch (e) {
      // this is handled in mutation hook
    }
  };

  return (
    <div className={styles.autocompleteContainer}>
      <div>
        <PrimereactAutocomplete
          inputClassName={classNames(styles.autocompleteInput, {
            [styles.notEditable]: !editable,
            [styles.disabled]: disabled,
            [styles.fitContent]: mode === 'fit-content',
          })}
          style={{ width: width ? `${width}px` : 'auto' }}
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

      {showCreateButton && (
        <div className={styles.createItemButton}>
          <button onClick={handleCreateItem}>
            Create {createItemType} &apos;{inputValue}&apos;
          </button>
        </div>
      )}
      <span
        ref={measureRef}
        className={classNames(styles.autocompleteInput)}
        style={{ top: -999999, position: 'absolute', width: 'auto' }}
      >
        {inputValue || placeholder}
      </span>
    </div>
  );
};

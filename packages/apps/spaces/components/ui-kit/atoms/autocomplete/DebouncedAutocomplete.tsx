import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import {
  AutoComplete as PrimereactAutocomplete,
  AutoCompleteChangeParams,
} from 'primereact/autocomplete';
import styles from './autocomplete.module.scss';
import classNames from 'classnames';
import { useDebouncedCallback } from 'use-debounce';
import { toast } from 'react-toastify';
import { DeleteIconButton } from '@spaces/atoms/icon-button';

export interface SuggestionItem {
  label: string;
  value: string;
}

interface CustomAutoCompleteProps {
  initialValue: string;
  createItemType?: string;
  suggestions: SuggestionItem[];
  onChange: (value: SuggestionItem) => void;
  onAddNew?: (item: SuggestionItem) => Promise<any>;
  onClearInput?: () => void;
  newItemLabel: string;
  editable?: boolean;
  disabled?: boolean;
  placeholder?: string;
  mode?: 'default' | 'fit-content' | 'invisible';
  onSearch: any;
  itemTemplate?: any;
}

export const DebouncedAutocomplete = ({
  initialValue,
  suggestions = [],
  onChange,
  onAddNew,
  onClearInput,
  createItemType = '',
  newItemLabel = '',
  editable,
  disabled,
  placeholder = '',
  mode = 'default',
  onSearch,
  itemTemplate,
  ...rest
}: CustomAutoCompleteProps) => {
  const [inputValue, setInputValue] = useState<string>(initialValue);
  const [width, setWidth] = useState<number>();
  const [showCreateButton, setShowCreateButton] = useState<boolean>(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);
  const handleInputChange = (event: AutoCompleteChangeParams) => {
    const newInputValue = event.value;
    setInputValue(newInputValue);
  };
  const debouncedSearch = useDebouncedCallback(
    // function
    (value) => {
      onSearch(value);
    },
    // delay in ms
    300,
  );

  useLayoutEffect(() => {
    if (mode === 'fit-content') {
      setWidth((measureRef?.current?.scrollWidth || 0) + 2);
    }
  }, [inputValue]);

  useEffect(() => {
    if (inputValue && editable && suggestions.length === 0) {
      setShowCreateButton(true);
    }
    if (
      suggestions.length ||
      !editable ||
      !inputValue.length ||
      inputValue === initialValue
    ) {
      setShowCreateButton(false);
    }
  }, [suggestions, inputValue, initialValue, editable]);

  useEffect(() => {
    console.log(
      'inputValue: ' + inputValue + '; initialValue: ' + initialValue,
    );
    if (inputValue !== initialValue && !editable) {
      setInputValue(initialValue);
    }
  }, [inputValue, initialValue, editable]);

  const handleSelectItem = (event: { value: SuggestionItem }) => {
    const selectedValue = event.value;
    console.log('🏷️ ----- selectedValue: ', selectedValue);
    setInputValue(selectedValue?.label);
    onChange(selectedValue);
  };

  const handleAddNew = () => {
    if (!onAddNew) {
      toast.error('New item could not be created');
      return;
    }
    const newItem = { label: inputValue, value: inputValue };
    onAddNew(newItem);
    setInputValue('');
    inputRef?.current?.focus();
  };

  const search = (event: any) => {
    debouncedSearch(event.query);
    setInputValue(event.query);
  };

  const handleCreateItem = async () => {
    if (!onAddNew) {
      toast.error('New item could not be created');
      return;
    }
    try {
      const newItem = await onAddNew({ value: inputValue, label: inputValue });
      if (newItem) {
        handleSelectItem({
          value: {
            label: newItem[newItemLabel],
            value: newItem.id,
          },
        });
        setInputValue(newItem[newItemLabel]);
        setShowCreateButton(false);
      }
    } catch (e) {
      // this is handled in mutation hook
    }
  };

  return (
    <div
      className={styles.autocompleteContainer}
      style={{ width: mode === 'invisible' ? '100%' : 'auto' }}
    >
      <div className={styles.autocompleteInputWrapper}>
        <PrimereactAutocomplete
          inputClassName={classNames(styles.autocompleteInput, {
            [styles.notEditable]: !editable,
            [styles.disabled]: disabled,
            [styles.fitContent]: mode === 'fit-content',
            [styles.invisible]: mode === 'invisible',
          })}
          style={{
            display: 'block',
            width: width
              ? `${width}px`
              : mode === 'invisible'
              ? '100%'
              : 'auto',
          }}
          disabled={!editable || disabled}
          value={inputValue}
          delay={300}
          placeholder={placeholder}
          suggestions={suggestions}
          onChange={handleInputChange}
          itemTemplate={(data) =>
            itemTemplate ? (
              itemTemplate(data)
            ) : (
              <span onClick={() => handleSelectItem(data)}>{data.label}</span>
            )
          }
          completeMethod={search}
          onSelect={handleSelectItem}
          onKeyUp={(event) => {
            if (
              event.key === 'Enter' &&
              !suggestions.find((item) => item.label === inputValue)
            ) {
              handleAddNew();
            }
          }}
          inputRef={inputRef}
        />

        {editable && inputValue && onClearInput && (
          <DeleteIconButton
            onDelete={() => {
              setInputValue('');
              onClearInput();
            }}
          />
        )}
      </div>

      {showCreateButton && onAddNew && (
        <div className={styles.createItemButton}>
          <button onClick={handleCreateItem}>
            Create {createItemType} &apos;{inputValue}&apos;
          </button>
        </div>
      )}
      <div
        ref={measureRef}
        className={classNames(styles.autocompleteInput, styles.measureInput)}
      >
        {inputValue || placeholder}
      </div>
    </div>
  );
};

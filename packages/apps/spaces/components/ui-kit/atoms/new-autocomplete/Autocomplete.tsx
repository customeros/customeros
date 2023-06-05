import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import styles from './autocomplete.module.scss';
import classNames from 'classnames';
import { useDebouncedCallback } from 'use-debounce';
import { useDetectClickOutside } from '@spaces/hooks/useDetectClickOutside';
import { DebouncedInput } from '@spaces/atoms/input';
import { AutocompleteSuggestionList } from '@spaces/atoms/new-autocomplete/autocomplete-suggestion-list';
import {useTimeout} from "primereact";

export interface SuggestionItem {
  label: string;
  value: string;
}

interface AutocompleteProps {
  initialValue: string;
  suggestions: SuggestionItem[];
  onChange: (value: SuggestionItem | undefined) => void;
  onDoubleClick?: () => void;
  onClearInput?: () => void;
  editable?: boolean;
  disabled?: boolean;
  placeholder?: string;
  mode?: 'fit-content' | 'full-width';
  loading: boolean;
  saving: boolean;
  onSearch: any;
  itemTemplate?: any;
}

export const Autocomplete = ({
  initialValue,
  suggestions = [],
  onChange,
  onDoubleClick,
  onClearInput,
  editable,
  disabled,
  placeholder = '',
  mode = 'fit-content',
  loading,
  saving,
  onSearch,
  itemTemplate,
  ...rest
}: AutocompleteProps) => {
  const [inputValue, setInputValue] = useState<string>(initialValue);
  const [width, setWidth] = useState<number>();
  const inputRef = useRef<HTMLInputElement>(null);
  const [openSuggestionList, setOpenSuggestionList] = useState(false);
  const [highlightedItemIndex, setHighlightedItemIndex] = useState<number>(0);
  const measureRef = useRef<HTMLDivElement>(null);
  const handleInputChange = (event: any) => {
    const newInputValue = event.target.value;
    setInputValue(newInputValue);
    if (newInputValue !== '') {
      debouncedSearch(newInputValue);
      setOpenSuggestionList(true);
    }
  };
  const debouncedSearch = useDebouncedCallback(
    // function
    (value) => {
      onSearch(value);
    },
    // delay in ms
    150,
  );

  useLayoutEffect(() => {
    if (mode === 'fit-content') {
      setWidth((measureRef?.current?.scrollWidth || 0) + 2);
    }
  }, [inputValue]);

  useEffect(() => {
    setInputValue(initialValue);
  }, [initialValue]);

  useEffect(() => {
    if (!inputValue.length) {
      setOpenSuggestionList(false);
    }
  }, [inputValue]);

  const handleSelectItem = (value: SuggestionItem | undefined) => {
    setInputValue(value?.label ?? '');
    setOpenSuggestionList(false);
    onChange(value);
  };

  useEffect(() => {
    setTimeout(() => {
      inputRef?.current?.focus();
      inputRef?.current?.setSelectionRange(0, inputRef?.current?.value.length);
    }, 0);
  }, []);

  const handleSetCursorAtTheEndOfInput = () => {
    const inputLength = inputRef?.current?.value.length || 0;
    inputRef?.current?.setSelectionRange(inputLength, inputLength);
  }


  const autocompleteWrapperRef = useRef<HTMLDivElement>(null);

  useDetectClickOutside(autocompleteWrapperRef, () => {
    setOpenSuggestionList(false);
    onDoubleClick && onDoubleClick();
    if (inputValue !== initialValue) {
      setInputValue(initialValue);

      if (inputValue.length === 0) {
        onClearInput && onClearInput();
      }
    }
  });

  const handleKeyDown = (event: any) => {
    const { key, currentTarget } = event;

    switch (key) {
      case 'Enter':
        handleSelectItem(suggestions[highlightedItemIndex]);
        break;
      case 'ArrowDown':
        handleSelectNextSuggestion({
          currentIndex: highlightedItemIndex,
          onIndexSelect: setHighlightedItemIndex,
        });

        break;
      case 'ArrowUp':
        handleSelectPrevSuggestion({
          currentIndex: highlightedItemIndex,
          onIndexSelect: setHighlightedItemIndex,
        });

        break;
      case 'Escape':
        setOpenSuggestionList(false);
        break;
    }
  };

  const handleSelectNextSuggestion = ({
    currentIndex,
    onIndexSelect,
  }: {
    currentIndex: number;
    onIndexSelect: (index: number) => void;
  }) => {
    let nextIndex;
    // select first item from the list -> if nothing is selected yet and there are available options
    if (currentIndex === -1 && suggestions?.length >= 0) {
      nextIndex = 0;
    }
    // select next item if currently selected item is not last on the list
    else if (suggestions.length - 1 > currentIndex) {
      nextIndex = currentIndex + 1;
    } else {
      nextIndex = suggestions.length - 1;
    }
    onIndexSelect(nextIndex);
    setInputValue(suggestions[nextIndex].label || '');
  };

  const handleSelectPrevSuggestion = ({ currentIndex, onIndexSelect }: any) => {
    // deselect list -> move focus back to input / previous context
    if (currentIndex === 0) {
      onIndexSelect(-1);
      setInputValue('');
      return -1;
    }
    // select prev
    if (currentIndex > 0) {
      onIndexSelect(currentIndex - 1);
      setInputValue(suggestions[currentIndex - 1]?.label || '');
    }
    setTimeout(() => {
      handleSetCursorAtTheEndOfInput()
    },0)
  };

  return (
    <div
      ref={autocompleteWrapperRef}
      className={styles.autocompleteContainer}
      style={{ width: mode === 'full-width' ? '100%' : 'auto' }}
    >
      <div className={styles.autocompleteInputWrapper}>
        <DebouncedInput
          {...rest}
          inputRef={inputRef}
          className={classNames(styles.autocompleteInput, {
            [styles.notEditable]: !editable,
            [styles.disabled]: disabled,
            [styles.fitContent]: mode === 'fit-content',
            [styles.fullWidth]: mode === 'full-width',
          })}
          customStyles={{
            display: 'block',
            width: width
              ? `${width}px`
              : mode === 'full-width'
              ? '100%'
              : 'auto',
          }}
          minLength={1}
          saving={saving}
          debounceTimeout={300}
          disabled={!editable || disabled}
          value={inputValue}
          placeholder={placeholder}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
        />

        <AutocompleteSuggestionList
          onSearchResultSelect={handleSelectItem}
          loadingSuggestions={loading}
          suggestions={suggestions}
          openSugestionList={openSuggestionList}
          selectedIndex={highlightedItemIndex}
          showEmpty={false}
          // onIndexChanged={(index: number | null) => {
          //   if (index === null) {
          //     inputRef?.current?.focus();
          //     setTimeout(() => {
          //       const cursorPosition = inputRef?.current?.value
          //         .length as number;
          //       inputRef?.current?.setSelectionRange(
          //         cursorPosition,
          //         cursorPosition,
          //       );
          //     }, 0);
          //   }
          // }}
        />
      </div>

      <div
        ref={measureRef}
        className={classNames(styles.autocompleteInput, styles.measureInput)}
      >
        {inputValue || placeholder}
      </div>
    </div>
  );
};

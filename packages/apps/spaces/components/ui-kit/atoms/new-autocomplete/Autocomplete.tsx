import React, {
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
  useCallback,
} from 'react';
import classNames from 'classnames';
import { useDebouncedCallback } from 'use-debounce';
import { useDetectClickOutside } from '@spaces/hooks/useDetectClickOutside';
import { DebouncedInput } from '@spaces/atoms/input';
import { AutocompleteSuggestionList } from '@spaces/atoms/new-autocomplete/autocomplete-suggestion-list';
import styles from './autocomplete.module.scss';
export interface SuggestionItem {
  label: string;
  value: string;
}

type KeyActions = {
  [key: string]: () => void;
};
interface AutocompleteProps {
  initialValue: string;
  suggestionsMatch: SuggestionItem[];
  suggestionsFuzzyMatch: SuggestionItem[];
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
  suggestionsMatch = [],
  suggestionsFuzzyMatch = [],
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
  const [width, setWidth] = useState<number | undefined>();
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);
  const autocompleteWrapperRef = useRef<HTMLDivElement>(null);
  const [openSuggestionList, setOpenSuggestionList] = useState(false);
  const [highlightedItemIndex, setHighlightedItemIndex] = useState<number>(0);

  const debouncedSearch = useDebouncedCallback((value: string) => {
    onSearch(value);
  }, 150);

  const handleInputChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const newInputValue = event.target.value;
      setInputValue(newInputValue);

      if (newInputValue !== '') {
        debouncedSearch(newInputValue);
        setOpenSuggestionList(true);
      } else {
        setOpenSuggestionList(false);
      }
    },
    [suggestionsMatch, debouncedSearch],
  );

  useLayoutEffect(() => {
    if (mode === 'fit-content') {
      setWidth((measureRef?.current?.scrollWidth || 0) + 2);
    }
  }, [inputValue, mode]);

  useEffect(() => {
    setInputValue(initialValue);
  }, [initialValue]);

  useEffect(() => {
    if (!inputValue.length) {
      setOpenSuggestionList(false);
    }
  }, [inputValue]);

  useEffect(() => {
    setTimeout(() => {
      inputRef?.current?.focus();
      inputRef?.current?.setSelectionRange(0, inputRef?.current?.value.length);
    }, 0);
  }, []);

  const handleSetCursorAtTheEndOfInput = useCallback(() => {
    const inputLength = inputRef?.current?.value.length || 0;
    inputRef?.current?.setSelectionRange(inputLength, inputLength);
  }, []);

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

  const handleSelectItem = useCallback(() => {
    const suggestions = suggestionsMatch.length
      ? suggestionsMatch
      : suggestionsFuzzyMatch;
    const selecteditem = suggestions[highlightedItemIndex];
    setInputValue(selecteditem?.label ?? '');
    setOpenSuggestionList(false);
    onChange(selecteditem);
  }, [onChange, suggestionsMatch, suggestionsFuzzyMatch, highlightedItemIndex]);

  useEffect(() => {
    setHighlightedItemIndex(0);
  }, [suggestionsMatch, suggestionsFuzzyMatch]);

  const handleSelectNextSuggestion = useCallback(() => {
    const suggestions = suggestionsMatch.length
      ? suggestionsMatch
      : suggestionsFuzzyMatch;

    if (!suggestions.length) return;

    setHighlightedItemIndex((currentIndex) => {
      let nextIndex;

      if (currentIndex === -1 && suggestions.length >= 0) {
        nextIndex = 0;
      } else if (suggestions.length - 1 > currentIndex) {
        nextIndex = currentIndex + 1;
      } else {
        nextIndex = suggestions.length - 1;
      }

      setInputValue(suggestions[nextIndex].label || '');
      return nextIndex;
    });
  }, [suggestionsMatch, suggestionsFuzzyMatch]);

  const handleSelectPrevSuggestion = useCallback(() => {
    const suggestions = suggestionsMatch.length
      ? suggestionsMatch
      : suggestionsFuzzyMatch;
    if (!suggestions.length) return;

    setHighlightedItemIndex((currentIndex) => {
      if (currentIndex === 0) {
        setInputValue('');
        return -1;
      }

      if (currentIndex > 0) {
        const prevIndex = currentIndex - 1;
        setInputValue(suggestions[prevIndex]?.label || '');
        setTimeout(handleSetCursorAtTheEndOfInput, 0);
        return prevIndex;
      }

      return currentIndex;
    });
  }, [suggestionsMatch, suggestionsFuzzyMatch]);

  const handleKeyDown = useCallback(
    (event: React.KeyboardEvent<HTMLInputElement>) => {
      const { key } = event;

      const keyActions: KeyActions = {
        Enter: () => handleSelectItem(),
        ArrowDown: () => handleSelectNextSuggestion(),
        ArrowUp: () => handleSelectPrevSuggestion(),
        Escape: () => {
          setOpenSuggestionList(false);
          setInputValue(initialValue);
        },
      };
      const action = keyActions[key];
      if (action) {
        action();
      }
    },
    [
      suggestionsMatch,
      suggestionsFuzzyMatch,
      highlightedItemIndex,
      handleSelectItem,
      handleSelectNextSuggestion,
      handleSelectPrevSuggestion,
      initialValue,
    ],
  );

  return (
    <div
      ref={autocompleteWrapperRef}
      className={classNames(styles.autocompleteContainer, {
        [styles.editable]: editable,
      })}
      style={{ width: mode === 'full-width' ? '100%' : 'auto' }}
    >
      <div className={styles.autocompleteInputWrapper}>
        <DebouncedInput
          {...rest}
          inputRef={inputRef}
          inlineMode
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
          searchTerm={inputValue}
          suggestionsMatch={suggestionsMatch}
          suggestionsFuzzyMatch={suggestionsFuzzyMatch}
          openSugestionList={openSuggestionList}
          selectedIndex={highlightedItemIndex}
        />
      </div>

      <div
        ref={measureRef}
        className={classNames(styles.autocompleteInput, styles.measureInput)}
        style={{
          width: mode === 'fit-content' ? 'auto' : '100%',
        }}
      >
        {inputValue}
      </div>
    </div>
  );
};

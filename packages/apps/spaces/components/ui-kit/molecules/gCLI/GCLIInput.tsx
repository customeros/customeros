import React, {
  KeyboardEventHandler,
  useEffect,
  useRef,
  useState,
} from 'react';

import { SuggestionList } from './suggestion-list';
import { useGCLI } from './context/GCLIContext';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch } from '@fortawesome/free-solid-svg-icons';
import { DebouncedInput } from '@spaces/atoms/input/DebouncedInput';
import styles from './GCLIInput.module.scss';
import { uuid4 } from '@sentry/utils';

// TODO
// Filtering:
// 1. executing filter action on enter is liniting next filter options to specific context eg. organisation
// 2. executing filter action on comma is grouping objects, no action executed (allows to execute action on all of those objects)
// objects should be initially of the same type

// TODO
//   add result context change mechanism

export const GCLIInput = () => {
  // TODO action simulation to be removed!

  const {
    label,
    icon,
    inputPlaceholder,
    existingTerms,
    loadSuggestions,
    loadingSuggestions,
    suggestionsLoaded,
    onItemsChange,
    selectedTermFormat,
    highlightedItemIndex,
    onHighlightedItemChange,
    onSelectNextSuggestion,
    onSelectPrevSuggestion,
  } = useGCLI();

  // todo use input value to create fill in effect on navigate through results by keyboard ??? do we even need that? is this a proper use case
  const [selectedValues, setSelectedValues] = useState(existingTerms ?? []);
  const [locationTerms, setLocationTerms] = useState([] as any[]);
  const [searchTerms, setSearchTerms] = useState([] as any[]);

  useEffect(() => {
    setLocationTerms(
      selectedValues.filter((item: any) => item.type === 'STATE'),
    );
    setSearchTerms(
      selectedValues.filter((item: any) => item.type === 'GENERIC'),
    );
  }, [selectedValues]);

  const [searchQuery, setSearchQuery] = useState('');

  const [suggestions, setSuggestions] = useState<Array<any>>([]);
  const [displayAction, setDisplayActions] = useState(false);
  const [selectedAction, setSelectedAction] = useState(-1);

  const inputRef = useRef<HTMLInputElement>(null);

  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = React.useRef(null);

  useEffect(() => {
    if (!loadingSuggestions && suggestionsLoaded) {
      setSuggestions(suggestionsLoaded);
    }
  }, [loadingSuggestions, suggestionsLoaded]);

  // HANDLERS FOR GENERAL ACTIONS
  const handleSearchResultSelect = (item: any, defaultAction: string) => {
    console.log('bunica');
    setDropdownOpen(false);
    const items = [...selectedValues, item];
    setSelectedValues(items);
    onItemsChange(items);
    setSearchQuery('');
    setSuggestions([]);
    inputRef?.current?.focus();
  };

  // HANDLERS FOR GENERAL ACTIONS
  const handleAsSimpleSearch = () => {
    if (!searchQuery) return;
    handleSearchResultSelect(
      {
        id: uuid4(),
        type: 'GENERIC',
        display: searchQuery,
      },
      'FILTER',
    );
  };
  // END HANDLERS FOR GENERAL ACTIONS

  const handleInputKeyDown = (event: any) => {
    const { key, currentTarget, target } = event;
    switch (key) {
      case 'Enter':
        handleAsSimpleSearch();
        break;
      case ',':
        // handleCreateGroup(event)
        break;
      case 'Backspace':
        // if(!!selectedOptions.length && target?.value?.length === 0) {
        //     const newSelected = selectedOptions.slice(0, -1)
        //     setSelectedOptions(newSelected)
        // }
        break;
      case 'ArrowDown':
        onSelectNextSuggestion({
          suggestions,
          currentIndex: highlightedItemIndex,
          onIndexSelect: onHighlightedItemChange,
        });
        break;
      case 'ArrowUp':
        // code for ArrowUp action
        break;
      case 'ArrowRight':
        // code for ArrowRight action
        break;
      case 'Escape':
        setDropdownOpen(false);
    }
  };
  const handleSearchResultsKeyDown = (event: KeyboardEvent, option: string) => {
    const { key, currentTarget } = event;
    switch (key) {
      case 'Enter':
        // // execute action
        // console.log('Enter action selection')
        // // todo manage input state o
        // setSelectedOptions([...selectedOptions, option])
        break;
      case 'ArrowDown':
        onSelectNextSuggestion({
          suggestions,
          currentIndex: highlightedItemIndex,
          onIndexSelect: onHighlightedItemChange,
        });
        break;
      case 'ArrowUp':
        onSelectPrevSuggestion({
          currentIndex: highlightedItemIndex,
          onIndexSelect: onHighlightedItemChange,
        });
        break;
      case 'ArrowRight':
        console.log('Arrow Right action');
        setDisplayActions(true);
        setSelectedAction(0);
        break;
      case 'ArrowLeft':
        // close action dropdown, back to results
        console.log('Arrow right action');
        setDisplayActions(false);
        break;
      case 'Escape':
        setDropdownOpen(false);
        break;
    }
    console.groupEnd();
  };
  const handleActionKeyDown: KeyboardEventHandler = ({ key }) => {
    console.log(key);
    switch (key) {
      case 'Enter':
        // execute action
        console.log('Enter action');
        break;
      case 'ArrowDown':
        // select next
        console.log('Arrow down action');
        setSelectedAction(selectedAction + 1);
        break;
      case 'ArrowUp':
        // select prev
        console.log('Arrow up action');
        setSelectedAction(selectedAction - 1);
        break;
      case 'ArrowLeft':
        // close action dropdown, back to results
        console.log('Arrow left on action list');
        setDisplayActions(false);
        setSelectedAction(-1);
        break;

        break;
      default:
        break;
    }
  };

  const handleInputChange = (event: any) => {
    if (!event.target.value) {
      return;
    }
    setSearchQuery(event.target.value);
    inputRef.current?.focus();
    setDropdownOpen(true);

    if (!event.target.value) {
      setSuggestions([]);
    } else {
      loadSuggestions(event.target.value);
    }
  };

  return (
    <div className={styles.gcli_wrapper}>
      <div className={styles.input_wrapper}>
        <div className={styles.input_label_icon}>
          {icon && <div className={styles.input_icon}>{icon}</div>}

          <div className={styles.input_label}>{label}</div>
        </div>

        <div className={styles.selected_terms_wrapper}>
          {locationTerms.length > 0 && <span>in:&nbsp;</span>}
          {locationTerms.map((e, index) => {
            return (
              <div
                className={styles.selected_term}
                key={index}
                onClick={() => {
                  const filters = selectedValues.filter(
                    (term, i) =>
                      term.type !== e.type ||
                      (term.type === e.type && term.id !== e.id),
                  );
                  setSelectedValues(filters);
                  onItemsChange(filters);
                  inputRef?.current?.focus();
                }}
              >
                <div className={styles.selected_term_text}>
                  {selectedTermFormat ? selectedTermFormat(e) : e.display}
                  {index < locationTerms.length - 1 ? ',' : ''}
                </div>
              </div>
            );
          })}
          {searchTerms.length > 0 && <span>contains:&nbsp;</span>}
          {searchTerms.map((e, index) => (
            <div
              className={styles.selected_term}
              key={index}
              onClick={() => {
                const filters = selectedValues.filter(
                  (term, i) =>
                    term.type !== e.type ||
                    (term.type === e.type && term.id !== e.id),
                );
                setSelectedValues(filters);
                onItemsChange(filters);
                inputRef?.current?.focus();
              }}
            >
              <div className={styles.selected_term_text}>
                {selectedTermFormat ? selectedTermFormat(e) : e.display}
                {index < searchTerms.length - 1 ? ',' : ''}
              </div>
            </div>
          ))}
        </div>

        <DebouncedInput
          placeholder={inputPlaceholder}
          className={styles.gcli_input_search}
          minLength={1}
          value={searchQuery}
          onChange={handleInputChange}
          onKeyDown={handleInputKeyDown}
          debounceTimeout={50}
        />

        <div className={styles.input_actions}>
          {!loadingSuggestions && searchQuery !== '' && (
            <button
              className={styles.search_button}
              onClick={handleAsSimpleSearch}
            >
              <FontAwesomeIcon
                icon={faSearch}
                style={{ marginRight: '10px' }}
              />{' '}
              Search
            </button>
          )}
        </div>
      </div>
      {/* END SELECTED OPTIONS */}

      {dropdownOpen && searchQuery !== '' && (
        <SuggestionList
          onSearchResultSelect={handleSearchResultSelect}
          onSearchResultsKeyDown={handleSearchResultsKeyDown}
          onActionKeyDown={handleActionKeyDown}
          loadingSuggestions={loadingSuggestions}
          suggestions={suggestions}
          selectedAction={selectedAction}
          displayAction={displayAction}
          highlightedIndex={highlightedItemIndex}
        />
      )}
    </div>
  );
};

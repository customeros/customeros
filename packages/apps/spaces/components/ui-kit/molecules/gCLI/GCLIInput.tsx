import React, { useEffect, useRef, useState } from 'react';

// import { SuggestionList } from './suggestion-list/SuggestionList';
import { useGCLI } from './context/GCLIContext';
import Search from '@spaces/atoms/icons/Search';
import { DebouncedInput } from '@spaces/atoms/input/DebouncedInput';
import styles from './GCLIInput.module.scss';
import { uuid4 } from '@sentry/utils';
import classNames from 'classnames';
import { SuggestionType } from '@spaces/molecules/gCLI/suggestion-list/types';

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
  } = useGCLI();

  // todo use input value to create fill in effect on navigate through results by keyboard ??? do we even need that? is this a proper use case
  const [selectedValues, setSelectedValues] = useState(existingTerms ?? []);
  const [locationTerms, setLocationTerms] = useState([] as any[]);
  const [searchTerms, setSearchTerms] = useState([] as any[]);
  const [suggestionListKeyDown, setSuggestionListKeyDown] =
    useState<boolean>(false);

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

  const inputRef = useRef<HTMLInputElement>(null);

  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = React.useRef(null);

  useEffect(() => {
    if (!loadingSuggestions && suggestionsLoaded) {
      setSuggestions(
        suggestionsLoaded.map((item: any) => {
          return {
            id: item.id,
            type: item.type,
            display: item.display,
            data: item.data,
            highlighted: false,
          } as SuggestionType;
        }),
      );
    }
  }, [loadingSuggestions, suggestionsLoaded]);

  // HANDLERS FOR GENERAL ACTIONS
  const handleSearchResultSelect = (item: any | undefined) => {
    if (item === undefined) {
      setDropdownOpen(false);
      inputRef.current?.focus();
      return;
    }
    setDropdownOpen(false);
    const items = [...selectedValues, { ...item, ...{ highlighted: false } }];
    items.sort((a, b) => {
      if (a.type === 'STATE' && b.type === 'GENERIC') {
        return -1;
      }
      if (a.type === 'STATE' && b.type === 'STATE') {
        if (a.display < b.display) return -1;
        if (a.display > b.display) return 1;
        return 0;
      }
      if (a.type === 'GENERIC' && b.type === 'GENERIC') {
        if (a.display < b.display) return -1;
        if (a.display > b.display) return 1;
        return 0;
      }
      return 0;
    });
    setSelectedValues(items);
    onItemsChange(items);
    setSearchQuery('');
    setSuggestions([]);
    inputRef?.current?.focus();
  };

  // HANDLERS FOR GENERAL ACTIONS
  const handleAsSimpleSearch = () => {
    if (!searchQuery) return;
    handleSearchResultSelect({
      id: uuid4(),
      type: 'GENERIC',
      display: searchQuery,
    });
  };
  // END HANDLERS FOR GENERAL ACTIONS

  const [deleteTermsMode, setDeleteTermsMode] = useState(0);
  const [deleteTermIndex, setDeleteTermIndex] = useState(0);

  const handleInputKeyDown = (event: any) => {
    const { key, currentTarget, target } = event;
    switch (key) {
      case 'Backspace':
        if (target.selectionStart === 0 && selectedValues.length > 0) {
          if (deleteTermsMode === 0) {
            enterDeleteTermsMode();
            return;
          }

          //deleting term
          if (deleteTermsMode === 1) {
            const items = [
              ...selectedValues.slice(0, deleteTermIndex),
              ...selectedValues.slice(deleteTermIndex + 1),
            ];

            setSelectedValues(items);
            onItemsChange(items);

            let newDeleteTermIndex = -1;

            if (deleteTermIndex + 1 <= items.length) {
              newDeleteTermIndex = deleteTermIndex;
            } else if (deleteTermIndex - 1 > 0) {
              newDeleteTermIndex = deleteTermIndex - 1;
            } else if (deleteTermIndex - 1 == 0) {
              newDeleteTermIndex = 0;
            }

            hightlightDeleteTerm(newDeleteTermIndex);
            setDeleteTermIndex(newDeleteTermIndex);

            return;
          }
        }
        break;
      case 'Enter':
        handleAsSimpleSearch();
        exitDeleteTermsMode();
        break;
      case 'ArrowLeft':
        if (deleteTermsMode === 0 && event.target.selectionStart === 0) {
          enterDeleteTermsMode();
          return;
        }
        if (deleteTermsMode === 1) {
          event.preventDefault();
          if (deleteTermIndex - 1 < 0) return;
          setDeleteTermIndex(deleteTermIndex - 1);
          hightlightDeleteTerm(deleteTermIndex - 1);
          return;
        }
        break;
      case 'ArrowRight':
        if (deleteTermsMode === 1) {
          if (deleteTermIndex + 1 > selectedValues.length) {
            exitDeleteTermsMode();
            return;
          }
          setDeleteTermIndex(deleteTermIndex + 1);
          hightlightDeleteTerm(deleteTermIndex + 1);
          event.preventDefault();
        }
        break;
      case 'ArrowUp':
        break;
      case 'ArrowDown':
        setSuggestionListKeyDown(!suggestionListKeyDown);
        break;
      case 'Escape':
        setDropdownOpen(false);
        exitDeleteTermsMode();
        break;
      default:
        exitDeleteTermsMode();
    }
  };

  const hightlightDeleteTerm = (termIndex: number) => {
    setSelectedValues((prevState) => {
      return prevState.map((item, index) => {
        return {
          ...item,
          highlighted: index === termIndex,
        };
      });
    });
  };

  const enterDeleteTermsMode = () => {
    setDeleteTermsMode(1);
    setDeleteTermIndex(selectedValues.length - 1);
    hightlightDeleteTerm(selectedValues.length - 1);
  };

  const exitDeleteTermsMode = () => {
    setDeleteTermsMode(0);
    setDeleteTermIndex(-1);
    hightlightDeleteTerm(-1);
  };

  const handleInputChange = (event: any) => {
    exitDeleteTermsMode();
    if (!event.target.value) {
      setDropdownOpen(false);
      setSuggestions([]);
      return;
    }
    setSearchQuery(event.target.value);
    inputRef.current?.focus();
    setDropdownOpen(true);

    loadSuggestions(event.target.value);
  };

  return (
    <div className={styles.gcli_wrapper}>
      <div className={styles.input_wrapper}>
        <div className={styles.input_label_icon}>
          {icon && <div className={styles.input_icon}>{icon}</div>}

          <div className={styles.input_label}>{label}</div>
        </div>

        <div className={styles.selected_terms_wrapper}>
          {locationTerms.length > 0 && <span className={styles.gray}>in</span>}
          {locationTerms.map((e, index) => {
            return (
              <div
                className={classNames(styles.selected_term, {
                  [styles.selected_term_highlighted]: e.highlighted,
                })}
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
          {searchTerms.length > 0 && (
            <span className={styles.gray}>contains</span>
          )}
          {searchTerms.map((e, index) => (
            <div
              className={classNames(styles.selected_term, {
                [styles.selected_term_highlighted]: e.highlighted,
              })}
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
          inputRef={inputRef}
          placeholder={
            inputPlaceholder ?? selectedValues.length === 0
              ? 'Search and filter here...'
              : 'Add more filters here...'
          }
          className={styles.gcli_input_search}
          minLength={1}
          value={searchQuery}
          onChange={handleInputChange}
          onKeyDown={handleInputKeyDown}
          onClick={() => {
            exitDeleteTermsMode();
          }}
          debounceTimeout={50}
        />
        
        <div className={styles.input_actions}>
          {searchQuery !== '' && (
            <button
              className={styles.search_button}
              onClick={handleAsSimpleSearch}
            >
              <Search height={16} style={{ marginRight: '10px' }} />
              Search
            </button>
          )}
        </div>
      </div>

      <div className={styles.count_items}>
        {/* Total items: {totalItems} */} Total items: 100000
        </div>
      {/*/!* END SELECTED OPTIONS *!/*/}

      {/*{dropdownOpen && searchQuery !== '' && (*/}
      {/*  <SuggestionList*/}
      {/*    onSearchResultSelect={handleSearchResultSelect}*/}
      {/*    loadingSuggestions={loadingSuggestions}*/}
      {/*    suggestions={suggestions}*/}
      {/*    suggestionListKeyDown={suggestionListKeyDown}*/}
      {/*    onIndexChanged={(index: number | null) => {*/}
      {/*      if (index === null) {*/}
      {/*        inputRef?.current?.focus();*/}
      {/*        setTimeout(() => {*/}
      {/*          const cursorPosition = inputRef?.current?.value*/}
      {/*            .length as number;*/}
      {/*          inputRef?.current?.setSelectionRange(*/}
      {/*            cursorPosition,*/}
      {/*            cursorPosition,*/}
      {/*          );*/}
      {/*        }, 0);*/}
      {/*      }*/}
      {/*    }}*/}
      {/*  />*/}
      {/*)}*/}
    </div>
  );
};

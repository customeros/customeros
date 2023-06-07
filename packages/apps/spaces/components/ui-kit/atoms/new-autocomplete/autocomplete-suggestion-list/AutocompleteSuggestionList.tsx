import React from 'react';
import styles from '../autocomplete.module.scss';
import { AutocompleteSuggestion } from '@spaces/atoms/new-autocomplete/autocomplete-suggestion-list/AutocompleteSugestion';
import { SuggestionItem } from '@spaces/atoms/new-autocomplete/Autocomplete';

interface SuggestionListProps {
  openSugestionList: boolean;
  onSearchResultSelect: (item: SuggestionItem | undefined) => void;
  loadingSuggestions: boolean;
  selectedIndex: number | null;
  searchTerm: string;
  suggestionsMatch: SuggestionItem[];
  suggestionsFuzzyMatch: SuggestionItem[];
}

export const AutocompleteSuggestionList = ({
  openSugestionList,
  onSearchResultSelect,
  loadingSuggestions,
  selectedIndex,
  searchTerm,
  suggestionsMatch,
  suggestionsFuzzyMatch,
}: SuggestionListProps) => {
  return (
    <>
      {openSugestionList && (
        <div className={styles.result_list}>
          <div className={styles.list_search_results_wrapper}>
            {loadingSuggestions && (
              <div className={styles.loading}>
                <div className={styles.lds_dual_ring}></div>
              </div>
            )}

            {!loadingSuggestions &&
              suggestionsMatch.map((suggestion, i: number) => (
                <AutocompleteSuggestion
                  key={suggestion.value}
                  active={i === selectedIndex}
                  item={suggestion}
                  onClick={(e) => onSearchResultSelect(suggestion)}
                />
              ))}

            {!loadingSuggestions && suggestionsMatch.length === 0 && (
              <div className={styles.list_search_results_empty}>
                {/*<div>{searchTerm} not found</div>*/}
                did you mean…?
              </div>
            )}

            {!loadingSuggestions &&
              !suggestionsMatch.length &&
              suggestionsFuzzyMatch.map((suggestion, i: number) => (
                <AutocompleteSuggestion
                  key={suggestion.value}
                  active={i === selectedIndex}
                  item={suggestion}
                  onClick={(e) => onSearchResultSelect(suggestion)}
                />
              ))}
          </div>
        </div>
      )}
    </>
  );
};

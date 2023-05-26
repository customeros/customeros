import React, { useEffect, useState } from 'react';
import { type } from '../data';
import { Suggestion } from './Sugestion';
import { SuggestionType } from './types';
import styles from '../GCLIInput.module.scss';

interface SuggestionListProps {
  suggestionListKeyDown: boolean;
  onSearchResultSelect: (item: any | undefined) => void;
  loadingSuggestions: boolean;
  onIndexChanged: (currentIndex: number | null) => void;
  suggestions: Array<SuggestionType>;
}

export const SuggestionList = ({
  suggestionListKeyDown,
  onSearchResultSelect,
  loadingSuggestions,
  onIndexChanged,
  suggestions,
}: SuggestionListProps) => {
  const [highlightedItemIndex, setHighlightedItemIndex] = useState<
    null | number
  >(0);

  useEffect(() => {
    if (!loadingSuggestions) {
      setHighlightedItemIndex(null);
    }
  }, [loadingSuggestions]);

  useEffect(() => {
    if (!loadingSuggestions) {
      setHighlightedItemIndex(0);
    }
  }, [suggestionListKeyDown]);

  useEffect(() => {
    onIndexChanged(highlightedItemIndex);
  }, [highlightedItemIndex]);

  const handleSearchResultsKeyDown = (
    event: KeyboardEvent,
    optionId: string,
  ) => {
    const { key, currentTarget } = event;
    switch (key) {
      case 'Enter':
        onSearchResultSelect(suggestions.filter((p) => p.id === optionId)[0]);
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
        onSearchResultSelect(undefined);
        break;
    }
  };

  const handleSelectNextSuggestion = ({ currentIndex, onIndexSelect }: any) => {
    console.log('currentIndex: ', currentIndex);
    console.log('suggestions: ', suggestions.length);
    // select first item from the list -> if nothing is selected yet and there are available options
    if (currentIndex === null && suggestions?.length >= 0) {
      onIndexSelect(0);
    }
    // select next item if currently selected item is not last on the list
    else if (suggestions.length - 1 > currentIndex) {
      onIndexSelect(currentIndex + 1);
    } else {
      onIndexSelect(suggestions.length - 1);
    }
  };

  const handleSelectPrevSuggestion = ({ currentIndex, onIndexSelect }: any) => {
    // deselect list -> move focus back to input / previous context
    if (currentIndex === 0) {
      onIndexSelect(null);
    }
    // select prev
    if (currentIndex > 0) {
      onIndexSelect(currentIndex - 1);
    }
  };

  return (
    <section className={styles.result_list}>
      <div className={styles.list_search_results_wrapper}>
        {loadingSuggestions && (
          <div className={styles.loading}>
            <div className={styles.lds_dual_ring}></div>
          </div>
        )}
        {!loadingSuggestions && suggestions.length === 0 && (
          <div className={styles.list_search_results_empty}>
            No results found. Type Enter to search.
          </div>
        )}
        {!loadingSuggestions &&
          suggestions.map((suggestion, i: number) => (
            <React.Fragment key={suggestion.type + '_' + suggestion.id}>
              <Suggestion
                key={suggestion.id}
                active={i === highlightedItemIndex}
                item={suggestion}
                onClick={(e) => onSearchResultSelect(suggestion)}
                onKeyDown={(e: any) =>
                  handleSearchResultsKeyDown(e, suggestion.id)
                }
                defaultAction={
                  type.find((e) => e.name === suggestion.type)?.defaultAction
                }
              />
            </React.Fragment>
          ))}
      </div>
    </section>
  );
};

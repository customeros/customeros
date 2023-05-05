import React from 'react';
import { type } from '../data';
import { Suggestion } from './Sugestion';
import { SuggestionType } from './types';
import styles from '../GCLIInput.module.scss';

interface SuggestionListProps {
  onSearchResultSelect: any;
  onSearchResultsKeyDown: any;
  onActionKeyDown: any;
  loadingSuggestions: boolean;
  suggestions: Array<SuggestionType>;
  highlightedIndex: number | null;
  selectedAction: number;
  displayAction: boolean;
}

export const SuggestionList = ({
  onSearchResultSelect,
  onSearchResultsKeyDown,
  onActionKeyDown,
  selectedAction,
  highlightedIndex,
  loadingSuggestions,
  suggestions,
  displayAction,
}: SuggestionListProps) => {
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
              {(!displayAction || i === highlightedIndex) && (
                <Suggestion
                  key={suggestion.id}
                  active={i === highlightedIndex}
                  item={suggestion}
                  onClick={(e) => onSearchResultSelect(suggestion)}
                  onKeyDown={(e) => onSearchResultsKeyDown(e, suggestion.id)}
                  defaultAction={
                    type.find((e) => e.name === suggestion.type)?.defaultAction
                  }
                />
              )}
            </React.Fragment>
          ))}
      </div>
    </section>
  );
};

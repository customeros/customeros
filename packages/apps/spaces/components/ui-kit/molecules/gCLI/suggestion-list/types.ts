export type { GCLIInputMode } from '../types';
export type SuggestionType = {
  id: string;
  type: string;
  display: string;
  data: SuggestionTypeData[];
  highlighted: boolean;
};
export type SuggestionTypeData = {
  key: string;
  value: string;
  display: string;
};

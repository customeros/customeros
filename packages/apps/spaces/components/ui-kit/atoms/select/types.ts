export interface SelectOption<T = string> {
  value: T;
  label: string;
}

export enum SelectActionType {
  'OPEN',
  'CLOSE',
  'TOGGLE',
  'KEYDOWN',
  'BLUR',
  'CLICK',
  'DBLCLICK',
  'CHANGE',
  'SELECT',
  'MOUSEENTER',
}

export type SelectState = {
  value: string;
  selection: string;
  isOpen: boolean;
  isEditing: boolean;
  currentIndex: number;
  items: SelectOption[];
  defaultItems: SelectOption[];
};

export type SelectAction = {
  type: SelectActionType;
  payload?: unknown;
};

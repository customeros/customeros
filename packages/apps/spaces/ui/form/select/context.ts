import type {
  ChangeEventHandler,
  FocusEventHandler,
  KeyboardEventHandler,
  MouseEventHandler,
  RefObject,
} from 'react';
import { createContext } from 'react';

import { SelectState } from './types';
import { defaultState } from './reducer';

export const noop = () => undefined;

interface SelectContextMethods {
  state: SelectState;
  defaultValue?: string;
  inputRef: RefObject<HTMLSpanElement> | null;
  menuRef: RefObject<HTMLUListElement> | null;
  autofillValue: string;
  getInputProps: () => {
    onBlur: FocusEventHandler<HTMLInputElement>;
    onInput: ChangeEventHandler<HTMLSpanElement>;
    onKeyDown: KeyboardEventHandler<HTMLInputElement>;
    onDoubleClick: MouseEventHandler<HTMLInputElement>;
  };
  getMenuProps: ({ maxHeight }: { maxHeight: number }) => {
    ref: RefObject<HTMLUListElement> | null;
  };
  getMenuItemProps: (options: { value: string; index: number }) => {
    onClick: MouseEventHandler<HTMLLIElement>;
    ref: RefObject<HTMLLIElement> | null;
  };
  getWrapperProps: () => {
    onClick: MouseEventHandler<HTMLDivElement>;
    ref: RefObject<HTMLDivElement> | null;
  };
}

export const SelectContext = createContext<SelectContextMethods>({
  inputRef: null,
  menuRef: null,
  state: defaultState,
  autofillValue: '',
  getInputProps: () => ({
    onBlur: noop,
    onDoubleClick: noop,
    onInput: noop,
    onKeyDown: noop,
  }),
  getMenuProps: () => ({
    ref: null,
  }),
  getMenuItemProps: () => ({
    onClick: noop,
    ref: null,
  }),
  getWrapperProps: () => ({
    onClick: noop,
    ref: null,
  }),
});

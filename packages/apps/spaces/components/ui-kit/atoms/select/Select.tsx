import type {
  PropsWithChildren,
  ChangeEventHandler,
  KeyboardEventHandler,
  FocusEventHandler,
  MouseEventHandler,
} from 'react';
import { useEffect, useRef, useReducer } from 'react';
import { useDetectClickOutside } from '@spaces/hooks/useDetectClickOutside';

import { reducer, defaultState } from './reducer';
import { SelectActionType, SelectOption, SelectState } from './types';
import { SelectContext } from './context';

interface SelectProps<T = string> {
  defaultValue?: T extends string ? string : undefined;
  value?: T extends string ? string : undefined;
  options: SelectOption[];
  onChange?: (value: string) => void;
  onSelect?: (selection: T) => void;
}

type InputType = HTMLSpanElement | HTMLInputElement;

function placeCaretAtEnd(el: HTMLElement) {
  el.focus();
  const range = document.createRange();
  range.selectNodeContents(el);
  range.collapse(false);
  const sel = window.getSelection();
  sel?.removeAllRanges();
  sel?.addRange(range);
}

export const Select = <T = string,>({
  options = [],
  children,
  value,
  defaultValue,
  onChange,
  onSelect,
}: PropsWithChildren<SelectProps<T>>) => {
  const inputRef = useRef<InputType>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLUListElement>(null);

  const [state, dispatch] = useReducer(reducer, {
    ...defaultState,
    selection: value ? value : defaultValue ?? '',
    items: options,
    defaultItems: options,
    defaultSelection: value ? value : defaultValue ?? '',
  } as SelectState<string>);

  const autofillValue = (() => {
    if (!state.value) return '';
    const item = state.items?.[0];
    if (!item) return '';

    const label = item.label;
    const value = state.value;
    const index = label.toLowerCase().indexOf(value.toLowerCase());

    return label.slice(index + value.length);
  })();

  const getInputProps = () => {
    const onInput: ChangeEventHandler<InputType> = (e) => {
      dispatch({
        type: SelectActionType.CHANGE,
        payload: e.target.textContent,
      });
      onChange?.(e.target.textContent ?? '');
    };

    const onKeyDown: KeyboardEventHandler<InputType> = (e) => {
      dispatch({ type: SelectActionType.KEYDOWN, payload: e.key });
      if (e.key === 'Enter') {
        const selection = state.items?.[state.currentIndex]?.value ?? '';
        onSelect?.(selection as T);
      }
    };

    const onBlur: FocusEventHandler<InputType> = () => {
      dispatch({ type: SelectActionType.BLUR });
    };

    const onDoubleClick: MouseEventHandler<InputType> = () => {
      dispatch({ type: SelectActionType.DBLCLICK });
      setTimeout(() => inputRef.current?.focus(), 0);
    };

    return {
      onInput,
      onKeyDown,
      onBlur,
      onDoubleClick,
      'data-dropdown': 'input',
      ref: inputRef,
    };
  };

  const getMenuProps = ({ maxHeight }: { maxHeight: number }) => {
    const style = {
      marginTop: inputRef?.current?.offsetHeight ?? undefined,
      visibility: state.isOpen ? 'visible' : 'hidden',
      maxHeight,
    };

    return {
      ref: menuRef,
      style,
    };
  };

  const getMenuItemProps = ({ value }: { value: string; index: number }) => {
    const onClick: MouseEventHandler<HTMLLIElement> = (e) => {
      e.preventDefault();
      dispatch({ type: SelectActionType.SELECT, payload: value });
      onSelect?.(value as T);
      inputRef.current?.focus();
    };

    const onMouseEnter: MouseEventHandler<HTMLLIElement> = () => {
      const index = state.items.findIndex((item) => item.value === value);
      dispatch({ type: SelectActionType.MOUSEENTER, payload: index });
    };

    return {
      onClick,
      onMouseEnter,
      ref: null,
      'data-dropdown': 'menuitem',
    };
  };

  const getWrapperProps = () => {
    const onClick: MouseEventHandler<HTMLDivElement> = (e) => {
      const targetEl = (e.target as Element).getAttribute('data-dropdown');
      dispatch({ type: SelectActionType.CLICK, payload: targetEl });
    };

    return {
      onClick,
      ref: wrapperRef,
    };
  };

  useDetectClickOutside(wrapperRef, () => {
    dispatch({ type: SelectActionType.CLOSE });
  });

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.textContent = state.selection
        ? options.find((o) => o.value === state.selection)?.label ?? ''
        : state.value;
      placeCaretAtEnd(inputRef.current as HTMLElement);
    }
    if (state.selection) {
      dispatch({
        type: SelectActionType.SET_DEFAULT_SELECTION,
        payload: state.selection,
      });
    }
  }, [state.selection, state.value, options]);

  useEffect(() => {
    dispatch({ type: SelectActionType.SET_SELECTION, payload: value });
  }, [value]);

  return (
    <SelectContext.Provider
      value={{
        state,
        inputRef,
        defaultValue,
        menuRef,
        autofillValue,
        getInputProps,
        getMenuProps,
        getMenuItemProps,
        getWrapperProps,
      }}
    >
      {children}
    </SelectContext.Provider>
  );
};

// noinspection CommaExpressionJS

import type {
  ChangeEventHandler,
  FocusEventHandler,
  KeyboardEventHandler,
  MouseEventHandler,
  PropsWithChildren,
} from 'react';
import { useEffect, useReducer, useRef, useState } from 'react';
import { useDetectClickOutside } from '@shared/hooks/useDetectClickOutside';

import { defaultState, reducer } from '../reducer';
import { SelectActionType, SelectOption, SelectState } from '../types';
import { SelectContext } from '../context';

interface SelectProps<T = string> {
  defaultValue?: T extends string ? string : undefined;
  value?: T extends string ? string : undefined;
  options: SelectOption[];
  onChange?: (value: string) => void;
  onSelect?: (selection: T) => void;
  onCreateNewOption?: (selection: T) => void;
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
function selectNodeContents(el: HTMLElement) {
  const range = document.createRange();
  range.selectNodeContents(el);
  const sel = window.getSelection();
  sel?.removeAllRanges();
  sel?.addRange(range);
}

export const Select = <T extends string>({
  options = [],
  children,
  value,
  defaultValue,
  onChange,
  onSelect,
  onCreateNewOption,
}: PropsWithChildren<SelectProps<T>>) => {
  const inputRef = useRef<InputType>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLUListElement>(null);
  const [clickingOption, setClickingOption] = useState(false);

  const [state, dispatch] = useReducer(reducer, {
    ...defaultState,
    selection: value ? value : defaultValue ?? '',
    items: options,
    defaultItems: options,
    defaultSelection: value ? value : defaultValue ?? '',
    isCreating: false,
    canCreate: onCreateNewOption !== undefined,
  } as SelectState<T>);

  const autofillValue = (() => {
    if (!state.value) return '';
    const item = state.items?.[0];
    if (!item) return '';
    const value = state.value;
    const shouldAutofill = item.label
      .toLowerCase()
      .startsWith(value.trim().toLowerCase());

    if (!shouldAutofill) return '';
    const label = item.label;

    const index = label.toLowerCase().indexOf(value.trim());

    if (index === -1) return '';
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
        if (state.canCreate && state.value && state.items.length === 0) {
          onCreateNewOption?.(state.value as T);
          return;
        }

        const selection = state.items?.[state.currentIndex]?.value ?? '';
        dispatch({ type: SelectActionType.SELECT, payload: selection });
        onSelect?.(selection as T);
      }

      if (e.key === 'Backspace' && state.selection.length) {
        dispatch({ type: SelectActionType.SELECT, payload: '' });
        onSelect?.('' as T);
      }
    };

    const onBlur: FocusEventHandler<InputType> = () => {
      if (clickingOption) {
        setClickingOption(false);
        return;
      }
      dispatch({ type: SelectActionType.BLUR });
    };

    const onDoubleClick: MouseEventHandler<InputType> = () => {
      dispatch({ type: SelectActionType.DBLCLICK });

      setTimeout(() => {
        inputRef.current?.focus();
        selectNodeContents(inputRef.current as HTMLInputElement);
      }, 0);
    };

    const onClick: MouseEventHandler<InputType> = () => {
      if (!value) {
        dispatch({ type: SelectActionType.SET_EDITABLE });
        setTimeout(() => {
          inputRef.current?.focus();
        }, 0);
      }
    };

    return {
      onInput,
      onKeyDown,
      onBlur,
      onDoubleClick,
      onClick,
      'data-dropdown': 'input',
      ref: inputRef,
    };
  };

  const getMenuProps = ({ maxHeight }: { maxHeight: number }) => {
    const style = {
      marginTop: inputRef?.current?.offsetHeight
        ? inputRef?.current?.offsetHeight + 6
        : undefined,
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

      if (state.canCreate && value && state.items.length === 0) {
        onCreateNewOption?.(state.value as T);
        return;
      }
      dispatch({ type: SelectActionType.SELECT, payload: value });
      onSelect?.(value as T);
      inputRef.current?.focus();
    };

    const onMouseDown: MouseEventHandler<HTMLLIElement> = () => {
      setClickingOption(true);
    };
    const onMouseEnter: MouseEventHandler<HTMLLIElement> = () => {
      if (state.items.length === 0 && state.canCreate) {
        dispatch({ type: SelectActionType.MOUSEENTER, payload: 0 });
        return;
      }

      const index = state.items.findIndex((item) => item.value === value);
      dispatch({ type: SelectActionType.MOUSEENTER, payload: index });
    };

    return {
      onClick,
      onMouseEnter,
      onMouseDown,
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
    dispatch({ type: SelectActionType.SET_DEFAULT_ITEMS, payload: options });
  }, [options]);

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

// noinspection CommaExpressionJS

import type {
  FocusEventHandler,
  KeyboardEventHandler,
  MouseEventHandler,
  PropsWithChildren,
  RefObject,
} from 'react';
import {
  createContext,
  useContext,
  useEffect,
  useReducer,
  useRef,
  useState,
} from 'react';
import { useDetectClickOutside } from '@shared/hooks/useDetectClickOutside';

import { defaultState, reducer } from '../../reducer';
import { SelectActionType, SelectOption, SelectState } from '../../types';

export const noop = () => undefined;

interface SelectContextMethods {
  state: SelectState;
  defaultValue?: string;
  toggleButtonRef: RefObject<HTMLSpanElement> | null;
  menuRef: RefObject<HTMLUListElement> | null;
  getToggleButtonProps: () => {
    onBlur: FocusEventHandler<HTMLInputElement>;
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
  toggleButtonRef: null,
  menuRef: null,
  state: defaultState,
  getToggleButtonProps: () => ({
    onBlur: noop,
    onDoubleClick: noop,
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

interface SelectProps<T = string> {
  defaultValue?: T extends string ? string : undefined;
  value?: T extends string ? string : undefined;
  options: SelectOption[];
  onSelect?: (selection: T) => void;
}

type InputType = HTMLSpanElement | HTMLInputElement;

export const SingleSelect = <T extends string>({
  options = [],
  children,
  value,
  defaultValue,
  onSelect,
}: PropsWithChildren<SelectProps<T>>) => {
  const toggleButtonRef = useRef<HTMLDivElement>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLUListElement>(null);
  const [clickingOption, setClickingOption] = useState(false);

  const [state, dispatch] = useReducer(reducer, {
    ...defaultState,
    selection: value ? value : defaultValue ?? '',
    items: options,
    defaultItems: options,
    defaultSelection: value ? value : defaultValue ?? '',
  } as SelectState<T>);

  const getToggleButtonProps = () => {
    const onKeyDown: KeyboardEventHandler<InputType> = (e) => {
      e.preventDefault();
      dispatch({ type: SelectActionType.KEYDOWN, payload: e.key });

      if (e.key === 'Enter') {
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
    };

    const onClick: MouseEventHandler<InputType> = () => {
      if (!value) {
        setTimeout(() => {
          dispatch({ type: SelectActionType.SET_EDITABLE });
          dispatch({ type: SelectActionType.OPEN });
        }, 0);
      }
    };

    return {
      onKeyDown,
      onBlur,
      onDoubleClick,
      onClick,
      'data-dropdown': 'input',
      ref: toggleButtonRef,
    };
  };

  const getMenuProps = ({ maxHeight }: { maxHeight: number }) => {
    const style = {
      marginTop: toggleButtonRef?.current?.offsetHeight
        ? toggleButtonRef?.current?.offsetHeight + 6
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
      dispatch({ type: SelectActionType.SELECT, payload: value });
      onSelect?.(value as T);
      toggleButtonRef.current?.focus();
    };

    const onMouseDown: MouseEventHandler<HTMLLIElement> = () => {
      setClickingOption(true);
    };
    const onMouseEnter: MouseEventHandler<HTMLLIElement> = () => {
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
        toggleButtonRef,
        defaultValue,
        menuRef,
        getToggleButtonProps,
        getMenuProps,
        getMenuItemProps,
        getWrapperProps,
      }}
    >
      {children}
    </SelectContext.Provider>
  );
};

export const useSingleSelect = () => useContext(SelectContext);

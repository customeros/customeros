import type {
  ChangeEventHandler,
  PropsWithChildren,
  KeyboardEventHandler,
  FocusEventHandler,
  MouseEventHandler,
  RefObject,
} from 'react';
import {
  useRef,
  createContext,
  useContext,
  useReducer,
  useEffect,
} from 'react';
import classNames from 'classnames';
import { useDetectClickOutside } from '@spaces/hooks/useDetectClickOutside';

import styles from './organization-relationship.module.scss';

const noop = () => undefined;

enum DropdownActionType {
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
type DropdownState = {
  value: string;
  selection: string;
  isOpen: boolean;
  isEditing: boolean;
  currentIndex: number;
  items: DropdownMenuOption[];
  defaultItems: DropdownMenuOption[];
};
type DropdownAction = {
  type: DropdownActionType;
  payload?: unknown;
};

interface DropdownProps {
  defaultValue?: string;
  options: DropdownMenuOption[];
}

const defaultState: DropdownState = {
  value: '',
  selection: '',
  isOpen: false,
  isEditing: false,
  currentIndex: -1,
  items: [],
  defaultItems: [],
};

const keyEventReducer = (state: DropdownState, key: string) => {
  if (!state.isEditing) return state;

  switch (key) {
    case 'ArrowDown':
      if (state.currentIndex === state.items.length - 1)
        return { ...state, isOpen: true };

      return {
        ...state,
        isOpen: true,
        currentIndex: !state.isOpen
          ? state.currentIndex
          : state.currentIndex + 1,
      };
    case 'ArrowUp':
      if (!state.isOpen) return state;
      if (state.currentIndex === 0)
        return { ...state, isOpen: false, currentIndex: -1 };

      return {
        ...state,
        currentIndex: state.currentIndex - 1,
      };
    case 'Escape':
      if (!state.isOpen) return { ...state, isEditing: false };
      return { ...state, isOpen: false };
    case 'Enter': {
      const selection = state.items?.[state.currentIndex]?.value ?? '';

      return {
        ...state,
        value: '',
        items: [...state.defaultItems],
        isOpen: false,
        isEditing: false,
        selection,
      };
    }
    case 'Backspace': {
      // const items = !state.value ? [...state.defaultItems] : state.items;
      if (state.selection) return { ...state, selection: '' };
      return state;
    }
    default:
      return state;
  }
};

const reducer = (state: DropdownState, action: DropdownAction) => {
  switch (action.type) {
    case DropdownActionType.OPEN:
      return { ...state, isOpen: true };
    case DropdownActionType.CLOSE:
      return { ...state, isOpen: false, isEditing: false };
    case DropdownActionType.TOGGLE:
      return { ...state, isOpen: !state.isOpen };
    case DropdownActionType.KEYDOWN:
      return keyEventReducer(state, action?.payload as string);
    case DropdownActionType.BLUR: {
      if (state.selection) return state;
      if (!state.value) return state;

      const selection = state.items?.[0]?.value ?? '';
      return {
        ...state,
        selection,
        value: '',
        items: [...state.defaultItems],
        currentIndex: -1,
      };
    }
    case DropdownActionType.DBLCLICK:
      return { ...state, isEditing: true, isOpen: true };
    case DropdownActionType.CLICK:
      switch (action.payload) {
        case 'input':
          return state;
        case 'menuitem':
          return { ...state, isOpen: false, isEditing: false };
        default:
          return { ...state, isOpen: false, isEditing: false };
      }
    case DropdownActionType.CHANGE: {
      const value = state.selection
        ? (action?.payload as string)[0]
        : (action?.payload as string);

      const items = (() => {
        return value
          ? [...state.defaultItems].filter((item) =>
              item.label
                .toLowerCase()
                .includes((action?.payload as string).toLowerCase()),
            )
          : state.defaultItems;
      })();

      return {
        ...state,
        value,
        items,
        selection: '',
        isOpen: true,
        currentIndex: value ? 0 : state.currentIndex,
      };
    }
    case DropdownActionType.SELECT:
      return {
        ...state,
        selection: action?.payload as string,
        value: '',
        items: [...state.defaultItems],
      };
    case DropdownActionType.MOUSEENTER:
      return { ...state, currentIndex: action?.payload as number };
    default:
      return state;
  }
};

const DropdownContext = createContext<{
  state: DropdownState;
  defaultValue?: string;
  inputRef: RefObject<HTMLSpanElement> | null;
  menuRef: RefObject<HTMLUListElement> | null;
  inputProps: {
    onBlur: FocusEventHandler<HTMLInputElement>;
    onInput: ChangeEventHandler<HTMLSpanElement>;
    onKeyDown: KeyboardEventHandler<HTMLInputElement>;
    onDoubleClick: MouseEventHandler<HTMLInputElement>;
  };
  getMenuProps: () => {
    ref: RefObject<HTMLUListElement> | null;
  };
  getMenuItemProps: (options: { value: string; index: number }) => {
    onClick: MouseEventHandler<HTMLLIElement>;
    ref: RefObject<HTMLLIElement> | null;
  };
}>({
  inputRef: null,
  menuRef: null,
  state: defaultState,
  inputProps: {
    onBlur: noop,
    onDoubleClick: noop,
    onInput: noop,
    onKeyDown: noop,
  },
  getMenuProps: () => ({
    ref: null,
  }),
  getMenuItemProps: () => ({
    onClick: noop,
    ref: null,
  }),
});

const useDropdown = () => useContext(DropdownContext);

const relationshipOptions = [
  { value: 'CUSTOMER', label: 'Customer' },
  { value: 'PROSPECT', label: 'Prospect' },
  { value: 'PARTNER', label: 'Partner' },
  { value: 'VENDOR', label: 'Vendor' },
  { value: 'OTHER', label: 'Other' },
  { value: 'CUSTOMER2', label: 'Customer2' },
  { value: 'PROSPECT2', label: 'Prospect2' },
  { value: 'PARTNER3', label: 'Partner2' },
  { value: 'VENDOR4', label: 'Vendor2' },
  { value: 'OTHER5', label: 'Other2' },
];

const Dropdown = ({
  options,
  children,
  defaultValue,
}: PropsWithChildren<DropdownProps>) => {
  const inputRef = useRef<HTMLSpanElement>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLUListElement>(null);

  const [state, dispatch] = useReducer(reducer, {
    ...defaultState,
    selection: defaultValue ?? '',
    items: options ?? [],
    defaultItems: options ?? [],
  });

  const inputProps = (() => {
    const onInput: ChangeEventHandler<HTMLSpanElement> = (e) => {
      dispatch({
        type: DropdownActionType.CHANGE,
        payload: e.target.textContent,
      });
    };

    const onKeyDown: KeyboardEventHandler<HTMLSpanElement> = (e) => {
      dispatch({ type: DropdownActionType.KEYDOWN, payload: e.key });
    };

    const onBlur: FocusEventHandler<HTMLInputElement> = () => {
      dispatch({ type: DropdownActionType.BLUR });
    };

    const onDoubleClick: MouseEventHandler<HTMLInputElement> = () => {
      dispatch({ type: DropdownActionType.DBLCLICK });
      setTimeout(() => inputRef.current?.focus(), 0);
    };

    return {
      onInput,
      onKeyDown,
      onBlur,
      onDoubleClick,
    };
  })();

  const getMenuProps = () => {
    return {
      ref: menuRef,
    };
  };

  const getMenuItemProps = ({
    value,
    index,
  }: {
    value: string;
    index: number;
  }) => {
    const onClick: MouseEventHandler<HTMLLIElement> = (e) => {
      e.preventDefault();
      dispatch({ type: DropdownActionType.SELECT, payload: value });
      inputRef.current?.focus();
    };

    const onMouseEnter: MouseEventHandler<HTMLLIElement> = () => {
      const index = state.items.findIndex((item) => item.value === value);
      dispatch({ type: DropdownActionType.MOUSEENTER, payload: index });
    };

    return {
      onClick,
      onMouseEnter,
      ref: null,
    };
  };

  const onClick: MouseEventHandler<HTMLDivElement> = (e) => {
    const targetEl = (e.target as Element).getAttribute('data-dropdown');
    dispatch({ type: DropdownActionType.CLICK, payload: targetEl });
  };

  useDetectClickOutside(wrapperRef, () => {
    dispatch({ type: DropdownActionType.CLOSE });
  });

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.textContent = state.selection
        ? relationshipOptions.find((o) => o.value === state.selection)?.label ??
          ''
        : state.value;
    }
  }, [state.selection, state.value]);

  return (
    <DropdownContext.Provider
      value={{
        state,
        inputRef,
        defaultValue,
        inputProps,
        menuRef,
        getMenuProps,
        getMenuItemProps,
      }}
    >
      <div
        ref={wrapperRef}
        className={classNames(styles.dropdownWrapper)}
        onClick={onClick}
      >
        {children}
      </div>
    </DropdownContext.Provider>
  );
};

interface DropdownMenuOption<T = string> {
  value: T;
  label: string;
}

interface DropdownMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
}

const DropdownMenu = ({
  noOfVisibleItems = 7,
  itemSize = 25,
}: DropdownMenuProps) => {
  const { state, inputRef, getMenuProps, getMenuItemProps } = useDropdown();
  const maxMenuHeight = itemSize * noOfVisibleItems;

  return (
    <ul
      className={styles.dropdownMenu}
      style={{
        marginTop: inputRef?.current?.offsetHeight ?? undefined,
        maxHeight: maxMenuHeight,
        visibility: state.isOpen ? 'visible' : 'hidden',
      }}
      {...getMenuProps()}
    >
      {state.items.length ? (
        state.items.map(({ value, label }, index) => (
          <li
            key={value}
            data-dropdown='menuitem'
            className={classNames(styles.dropdownMenuItem, {
              [styles.isFocused]: state.currentIndex === index,
              [styles.isSelected]: state.selection === value,
            })}
            {...getMenuItemProps({ value, index })}
          >
            {label}
          </li>
        ))
      ) : (
        <li className={styles.dropdownMenuItem} data-dropdown='menuitem'>
          No options available
        </li>
      )}
    </ul>
  );
};

const DropdownInput = () => {
  const { state, inputRef, inputProps } = useDropdown();

  const autofillValue = (() => {
    if (!state.value) return '';
    const item = state.items?.[0];
    if (!item) return '';

    const label = item.label;
    const value = state.value;
    const index = label.toLowerCase().indexOf(value.toLowerCase());

    return label.slice(index + value.length);
  })();

  return (
    <>
      <span
        ref={inputRef}
        role='textbox'
        data-dropdown='input'
        placeholder='Relationship'
        contentEditable={state.isEditing}
        className={classNames(styles.dropdownInput)}
        {...inputProps}
      />
      <span className={styles.autofill}>{autofillValue}</span>
    </>
  );
};

interface OrganizationRelationshipProps {
  defaultValue?: string;
}

export const OrganizationRelationship = ({
  defaultValue,
}: OrganizationRelationshipProps) => {
  return (
    <Dropdown defaultValue={defaultValue} options={relationshipOptions}>
      <DropdownInput />
      <DropdownMenu />
    </Dropdown>
  );
};

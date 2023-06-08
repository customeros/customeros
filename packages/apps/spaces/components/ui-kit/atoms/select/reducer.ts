import { SelectState, SelectAction, SelectActionType } from './types';

export const defaultState: SelectState = {
  value: '',
  selection: '',
  isOpen: false,
  isEditing: false,
  currentIndex: -1,
  items: [],
  defaultItems: [],
  defaultSelection: '',
};

const keyEventReducer = (state: SelectState, key: string) => {
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
      return {
        ...state,
        isOpen: false,
        isEditing: false,
        selection:
          !state.value && !state.selection
            ? state.defaultSelection
            : state.selection,
      };
    case 'Enter': {
      const selection = !state.value
        ? ''
        : state.items?.[state.currentIndex]?.value ?? '';

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
      if (state.selection) return { ...state, value: '' };
      return state;
    }
    default:
      return state;
  }
};

export const reducer = (state: SelectState, action: SelectAction) => {
  switch (action.type) {
    case SelectActionType.OPEN:
      return { ...state, isOpen: true };
    case SelectActionType.CLOSE:
      return { ...state, isOpen: false, isEditing: false };
    case SelectActionType.TOGGLE:
      return { ...state, isOpen: !state.isOpen };
    case SelectActionType.KEYDOWN:
      return keyEventReducer(state, action?.payload as string);
    case SelectActionType.BLUR: {
      if (state.selection) return state;

      const selection = !state.value
        ? state.defaultSelection
        : state.items?.[0]?.value ?? '';

      return {
        ...state,
        selection,
        value: '',
        items: [...state.defaultItems],
        currentIndex: -1,
      };
    }
    case SelectActionType.DBLCLICK:
      return { ...state, isEditing: true, isOpen: true };
    case SelectActionType.CLICK:
      switch (action.payload) {
        case 'input':
          return state;
        case 'menuitem':
          return { ...state, isOpen: false, isEditing: false };
        default:
          return { ...state, isOpen: false, isEditing: false };
      }
    case SelectActionType.CHANGE: {
      const value = action?.payload as string;

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
    case SelectActionType.SELECT:
      return {
        ...state,
        selection: action?.payload as string,
        value: '',
        items: [...state.defaultItems],
      };
    case SelectActionType.MOUSEENTER:
      return { ...state, currentIndex: action?.payload as number };
    case SelectActionType.RESET:
      return {
        ...defaultState,
        defaultValue: action.payload as string,
        defaultItems: state.defaultItems,
        items: state.defaultItems,
      };
    case SelectActionType.SET_SELECTION:
      return {
        ...state,
        selection: action.payload as string,
      };
    case SelectActionType.SET_DEFAULT_SELECTION:
      return {
        ...state,
        defaultSelection: action.payload as string,
      };
    default:
      return state;
  }
};

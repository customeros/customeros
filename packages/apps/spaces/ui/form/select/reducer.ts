import { SelectState, SelectAction, SelectActionType } from './types';

export const defaultState: SelectState = {
  value: '',
  selection: '',
  isOpen: false,
  isEditing: false,
  isCreating: false,
  canCreate: false,
  currentIndex: -1,
  items: [],
  defaultItems: [],
  defaultSelection: '',
};

const keyEventReducer = (state: SelectState, key: string) => {
  if (!state.isEditing) return state;

  switch (key) {
    case 'ArrowDown':
      if (state.items.length === 0 && state.canCreate && state.isCreating) {
        return { ...state, isOpen: true, currentIndex: 0 };
      }
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
        isCreating: false,
        value: '',
        selection: !state.selection ? state.defaultSelection : state.selection,
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
      if (state.selection)
        return {
          ...state,
          value: '',
          selection: '',
          currentIndex: -1,
          items: state.defaultItems,
        };
      if (!state.value)
        return {
          ...state,
          selection: '',
          currentIndex: -1,
          items: state.defaultItems,
        };
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
      return {
        ...state,
        isOpen: false,
        isEditing: false,
        items: state.defaultItems,
      };
    case SelectActionType.TOGGLE:
      return { ...state, isOpen: !state.isOpen };
    case SelectActionType.KEYDOWN:
      return keyEventReducer(state, action?.payload as string);
    case SelectActionType.BLUR: {
      if (state.selection) return state;

      if (!state.selection?.length && !state.value?.length)
        return {
          ...state,
          items: [...state.defaultItems],
          currentIndex: -1,
        };

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
    case SelectActionType.SET_EDITABLE:
      return { ...state, isEditing: true, isOpen: false };
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
          ? [...state.defaultItems]
              .filter((item) => {
                return item.label
                  .toLowerCase()
                  .includes(value.trim().toLowerCase());
              })
              .sort((a, b) => {
                if (
                  a.label.toLowerCase().startsWith(value.trim().toLowerCase())
                )
                  return -1;
                if (
                  b.label.toLowerCase().startsWith(value.trim().toLowerCase())
                )
                  return 1;
                return 0;
              })
          : state.defaultItems;
      })();

      return {
        ...state,
        value,
        items,
        isCreating: state.canCreate && !!value.length,
        selection: '',
        isOpen: true,
        currentIndex: value || state.canCreate ? 0 : -1,
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
    case SelectActionType.SET_DEFAULT_ITEMS:
      return {
        ...state,
        items: [...(action.payload as Array<any>)],
        defaultItems: [...(action.payload as Array<any>)],
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

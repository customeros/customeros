import { CSSWithMultiValues } from '@ui/utils';
import { OptionProps } from '@ui/form/SyncSelect';

export const contractButtonSelect = {
  singleValue: (props: CSSWithMultiValues) => {
    return {
      ...props,
      maxHeight: '22px',
      p: 0,
      minH: 'auto',
      color: 'inherit',
    };
  },

  input: (props: CSSWithMultiValues) => {
    return {
      ...props,
      maxHeight: '22px',
      minH: 'auto',
      p: 0,
    };
  },
  inputContainer: (props: CSSWithMultiValues) => {
    return {
      ...props,
      maxHeight: '22px',
      minH: 'auto',
      p: 0,
    };
  },

  control: (props: CSSWithMultiValues) => {
    return {
      ...props,
      w: '100%',
      border: 'none',
    };
  },
  option: (props: CSSWithMultiValues, state: OptionProps) => {
    return {
      ...props,
      my: '2px',
      borderRadius: 'md',
      color: 'gray.500',
      noOfLines: 1,
      // '-webkit-box-align': 'start',
      webkitBoxAlign: 'start',
      bg: state.isSelected ? 'primary.50' : 'white',
      boxShadow: state.isFocused ? 'menuOptionsFocus' : 'none',
      _hover: { bg: state.isSelected ? 'primary.50' : 'gray.100' },
    };
  },
};

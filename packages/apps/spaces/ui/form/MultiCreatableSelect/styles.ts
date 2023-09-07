import omit from 'lodash/omit';
import { ChakraStylesConfig, GroupBase } from 'chakra-react-select';
export const multiCreatableSelectStyles = (
  chakraStyles:
    | ChakraStylesConfig<unknown, boolean, GroupBase<unknown>>
    | undefined,
) => ({
  // @ts-expect-error fixme
  multiValue: (base) => ({
    ...base,
    padding: 0,
    paddingLeft: 2,
    paddingRight: 2,
    gap: 0,
    color: 'gray.500',
    background: 'gray.100',
    border: '1px solid',
    borderColor: 'gray.200',
    fontSize: 'md',
    marginRight: 1,
    cursor: 'default',
  }),
  // @ts-expect-error fixme
  clearIndicator: (base) => ({
    ...base,
    background: 'transparent',
    color: 'transparent',
    display: 'none',
  }),
  // @ts-expect-error fixme
  multiValueRemove: (styles, { data }) => ({
    ...styles,
    // visibility: 'hidden',
  }),
  // @ts-expect-error fixme
  container: (props) => ({
    ...props,
    minWidth: '300px',
    width: '100%',
    overflow: 'visible',
    _focusVisible: { border: 'none !important' },
    _focus: { border: 'none !important' },
  }),
  // @ts-expect-error fixme
  menuList: (props) => ({
    ...props,
    padding: '2',
    boxShadow: 'md',
    borderColor: 'gray.200',
    borderRadius: 'lg',
    maxHeight: '12rem',
  }),
  // @ts-expect-error fixme
  option: (props, { isSelected, isFocused }) => ({
    ...props,
    my: '2px',
    borderRadius: 'md',
    color: 'gray.700',
    bg: isSelected ? 'primary.50' : 'white',
    boxShadow: isFocused ? 'menuOptionsFocus' : 'none',
    _hover: { bg: isSelected ? 'primary.50' : 'gray.100' },
  }),
  // @ts-expect-error fixme
  groupHeading: (props) => ({
    ...props,
    color: 'gray.400',
    textTransform: 'uppercase',
    fontWeight: 'regular',
  }),
  // @ts-expect-error fixme
  input: (props) => ({
    ...props,
    color: 'gray.500',
    fontWeight: 'regular',
  }),
  // @ts-expect-error fixme
  valueContainer: (props) => ({
    ...props,
    maxH: '86px',
    overflowY: 'auto',
  }),
  ...omit<ChakraStylesConfig<unknown, false, GroupBase<unknown>>>(
    chakraStyles,
    [
      'container',
      'multiValueRemove',
      'multiValue',
      'clearIndicator',
      'menuList',
      'option',
      'groupHeading',
      'input',
      'valueContainer',
    ],
  ),
});

import omit from 'lodash/omit';
import { ChakraStylesConfig, GroupBase } from 'chakra-react-select';
import { CSSWithMultiValues } from '@chakra-ui/react';
export const multiCreatableSelectStyles = (
  chakraStyles:
    | ChakraStylesConfig<unknown, boolean, GroupBase<unknown>>
    | undefined,
) => ({
  multiValue: (base: CSSWithMultiValues) => ({
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
  clearIndicator: (base: CSSWithMultiValues) => ({
    ...base,
    background: 'transparent',
    color: 'transparent',
    display: 'none',
  }),
  multiValueRemove: (styles: CSSWithMultiValues) => ({
    ...styles,
  }),
  container: (props: CSSWithMultiValues) => ({
    ...props,
    minWidth: '300px',
    width: '100%',
    overflow: 'visible',
    _focusVisible: { border: 'none !important' },
    _focus: { border: 'none !important' },
  }),
  menuList: (props: CSSWithMultiValues) => ({
    ...props,
    padding: '2',
    boxShadow: 'md',
    borderColor: 'gray.200',
    borderRadius: 'lg',
    maxHeight: '12rem',
  }),
  option: (
    props: CSSWithMultiValues,
    { isSelected, isFocused }: { isSelected: boolean; isFocused: boolean },
  ) => ({
    ...props,
    my: '2px',
    borderRadius: 'md',
    color: 'gray.700',
    bg: isSelected ? 'primary.50' : 'white',
    boxShadow: isFocused ? 'menuOptionsFocus' : 'none',
    _hover: { bg: isSelected ? 'primary.50' : 'gray.100' },
  }),
  groupHeading: (props: CSSWithMultiValues) => ({
    ...props,
    color: 'gray.400',
    textTransform: 'uppercase',
    fontWeight: 'regular',
  }),
  input: (props: CSSWithMultiValues) => ({
    ...props,
    color: 'gray.500',
    fontWeight: 'regular',
  }),
  valueContainer: (props: CSSWithMultiValues) => ({
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

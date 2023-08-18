import { forwardRef, useCallback, useMemo } from 'react';
import {
  Select as _Select,
  Props,
  ControlProps,
  SelectInstance,
  chakraComponents,
  ClearIndicatorProps,
  GroupBase,
  ChakraStylesConfig,
} from 'chakra-react-select';

import Delete from '@spaces/atoms/icons/Delete';
import omit from 'lodash/omit';

export interface SelectProps extends Props<any, any, any> {
  leftElement?: React.ReactNode;
}

export const Select = forwardRef<SelectInstance, SelectProps>(
  ({ leftElement, chakraStyles, ...props }, ref) => {
    const Control = useCallback(({ children, ...rest }: ControlProps) => {
      return (
        <chakraComponents.Control {...rest}>
          {leftElement}
          {children}
        </chakraComponents.Control>
      );
    }, []);
    const ClearIndicator = useCallback(
      ({ children, ...rest }: ClearIndicatorProps) => {
        if (!rest.isFocused) return null;
        return (
          <chakraComponents.ClearIndicator {...rest} className='clearButton'>
            <Delete color='var(--chakra-colors-gray-500)' height='1rem' />
          </chakraComponents.ClearIndicator>
        );
      },
      [],
    );

    const components = useMemo(
      () => ({
        Control,
        ClearIndicator,
        DropdownIndicator: () => null,
      }),
      [Control],
    );

    return (
      <_Select
        variant='flushed'
        ref={ref}
        components={components}
        tabSelectsValue={false}
        chakraStyles={{
          container: (props, state) => ({
            ...props,
            w: '100%',
            overflow: 'visible',
            _hover: { cursor: 'pointer' },
            ...chakraStyles?.container?.(props, state),
          }),
          clearIndicator: (props, state) => ({
            ...props,
            padding: 2,
            _hover: {
              bg: 'gray.100',
            },
            ...chakraStyles?.clearIndicator?.(props, state),
          }),
          placeholder: (props) => ({
            ...props,
            color: 'gray.400',
          }),
          menuList: (props, state) => ({
            ...props,
            padding: '2',
            boxShadow: 'md',
            borderColor: 'gray.200',
            borderRadius: 'lg',
            ...chakraStyles?.menuList?.(props, state),
          }),
          option: (props, state) => ({
            ...props,
            my: '2px',
            borderRadius: 'md',
            color: 'gray.700',
            bg: state.isSelected ? 'primary.50' : 'white',
            boxShadow: state.isFocused ? 'menuOptionsFocus' : 'none',
            _hover: { bg: state.isSelected ? 'primary.50' : 'gray.100' },
            ...chakraStyles?.option?.(props, state),
          }),
          multiValue: (props, state) => ({
            ...props,
            borderRadius: 'full',
            bg: 'gray.50',
            color: 'gray.500',
            ml: 0,
            mr: 1,
            border: '1px solid',
            borderColor: 'gray.200',
            ...chakraStyles?.multiValue?.(props, state),
          }),
          groupHeading: (props, state) => ({
            ...props,
            color: 'gray.400',
            textTransform: 'uppercase',
            fontWeight: 'regular',
            ...chakraStyles?.groupHeading?.(props, state),
          }),
          ...omit<ChakraStylesConfig<unknown, false, GroupBase<unknown>>>(
            chakraStyles,
            [
              'container',
              'clearIndicator',
              'menuList',
              'option',
              'groupHeading',
            ],
          ),
        }}
        {...props}
      />
    );
  },
);

export type { SelectInstance };

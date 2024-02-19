import { useMemo, forwardRef, useCallback } from 'react';

import omit from 'lodash/omit';
import {
  Props,
  GroupBase,
  ControlProps,
  SelectInstance,
  chakraComponents,
  Select as _Select,
  ChakraStylesConfig,
  ClearIndicatorProps,
  LoadingIndicatorProps,
} from 'chakra-react-select';

import { Icons } from '@ui/media/Icon';

// Exhaustively typing this Props interface does not offer any benefit at this moment
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface SelectProps extends Props<any, any, any> {
  leftElement?: React.ReactNode;
}

// NOTE: custom elements like Control or ClearIndicator are wrapped in a useCallback without any dependencies on purpose.
// This is to avoid re-renders which cause the internal state of react-select to send up in weird states.
// Examples: Adding leftElement in the dependency array of the Control component will make the select unable to un-focus after selecting an option.

// Ideally, custom components should be declared outside of the Select component as per documentation. https://react-select.com/components#defining-components

export const Select = forwardRef<SelectInstance, SelectProps>(
  ({ leftElement, chakraStyles, components: _components, ...props }, ref) => {
    const Control = useCallback(
      ({ children, ...rest }: ControlProps) => {
        return (
          <chakraComponents.Control {...rest}>
            {leftElement}
            {children}
          </chakraComponents.Control>
        );
      },
      [leftElement],
    );
    const ClearIndicator = useCallback(
      ({ children, ...rest }: ClearIndicatorProps) => {
        const boxSize = (() => {
          switch (rest.selectProps.size) {
            case 'sm':
              return '3';
            case 'md':
              return '4';
            case 'lg':
              return '5';
            default:
              return '4';
          }
        })();

        if (!rest.isFocused) return null;

        return (
          <chakraComponents.ClearIndicator {...rest} className='clearButton'>
            <Icons.Delete color='gray.500' boxSize={boxSize} />
          </chakraComponents.ClearIndicator>
        );
      },
      [],
    );
    const LoadingIndicator = useCallback((props: LoadingIndicatorProps) => {
      return <chakraComponents.LoadingIndicator thickness='1px' {...props} />;
    }, []);

    const components = useMemo(
      () => ({
        Control,
        ClearIndicator,
        LoadingIndicator,
        DropdownIndicator: () => null,
        ..._components,
      }),
      [Control, _components],
    );

    return (
      <_Select
        variant='flushed'
        ref={ref}
        components={components}
        tabSelectsValue={false}
        chakraStyles={{
          container: (props, state) => {
            const readOnlyStyles = state.selectProps?.isReadOnly
              ? { pointerEvents: 'none' }
              : {};

            return {
              ...props,
              w: '100%',
              overflow: 'visible',
              _hover: { cursor: 'pointer' },
              ...chakraStyles?.container?.(props, state),
              ...readOnlyStyles,
            };
          },
          clearIndicator: (props, state) => ({
            ...props,
            padding: 2,
            _hover: {
              bg: '#F2F4F7',
            },
            ...chakraStyles?.clearIndicator?.(
              {
                ...props,
                padding: 2,
                _hover: {
                  bg: '#F2F4F7',
                },
              },
              state,
            ),
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
            '&::-webkit-scrollbar': {
              width: '4px',
              height: '4px',
              background: 'transparent',
            },
            '&::-webkit-scrollbar-track': {
              width: '4px',
              height: '4px',
              background: 'transparent',
            },
            '&::-webkit-scrollbar-thumb': {
              background: 'gray.500',
              borderRadius: '8px',
            },
            ...chakraStyles?.menuList?.(
              {
                ...props,
                padding: '2',
                boxShadow: 'md',
                borderColor: 'gray.200',
                borderRadius: 'lg',
                '&::-webkit-scrollbar': {
                  width: '4px',
                  height: '4px',
                  background: 'transparent',
                },
                '&::-webkit-scrollbar-track': {
                  width: '4px',
                  height: '4px',
                  background: 'transparent',
                },
                '&::-webkit-scrollbar-thumb': {
                  background: 'gray.500',
                  borderRadius: '8px',
                },
              },
              state,
            ),
          }),
          option: (props, state) => {
            return {
              ...props,
              my: '2px',
              borderRadius: 'md',
              color: 'gray.700',
              noOfLines: 1,
              WebkitBoxAlign: 'start',
              bg: state.isSelected ? 'gray.50' : 'white',
              boxShadow: state.isFocused ? 'menuOptionsFocus' : 'none',
              _hover: { bg: state.isSelected ? 'gray.50' : 'gray.50' },
              _selected: {
                bg: 'gray.50',
                fontWeight: 'medium',
                color: 'gray.700',
              },
              ...chakraStyles?.option?.(props, state),
            };
          },
          multiValue: (props, state) => ({
            ...props,
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
          loadingIndicator: (props, state) => ({
            ...props,
            color: 'gray.500',
            ...chakraStyles?.loadingIndicator?.(props, state),
          }),
          ...omit<ChakraStylesConfig<unknown, false, GroupBase<unknown>>>(
            chakraStyles,
            [
              'container',
              'loadingIndicator',
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

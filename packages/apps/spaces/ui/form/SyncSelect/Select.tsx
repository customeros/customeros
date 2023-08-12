import { forwardRef, useCallback, useMemo } from 'react';
import {
  Select as _Select,
  Props,
  ControlProps,
  SelectInstance,
  chakraComponents,
  ClearIndicatorProps,
} from 'chakra-react-select';
import Delete from '@spaces/atoms/icons/Delete';

export interface SelectProps extends Props<any, any, any> {
  leftElement?: React.ReactNode;
}

export const Select = forwardRef<SelectInstance, SelectProps>(
  ({ leftElement, ...props }, ref) => {
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
          container: (props) => ({
            ...props,
            w: '100%',
            overflow: 'visible',
            _hover: { cursor: 'pointer' },
          }),
          clearIndicator: (props) => ({
            ...props,
            padding: 2,
            _hover: {
              bg: 'gray.100',
            },
          }),
          menuList: (props) => ({
            ...props,
            padding: '2',
            boxShadow: 'md',
            borderColor: 'gray.200',
            borderRadius: 'lg',
          }),
          option: (props, { isSelected, isFocused }) => ({
            ...props,
            my: '2px',
            borderRadius: 'md',
            color: 'gray.700',
            bg: isSelected ? 'primary.50' : 'white',
            boxShadow: isFocused ? 'menuOptionsFocus' : 'none',
            _hover: { bg: isSelected ? 'primary.50' : 'gray.100' },
          }),
          groupHeading: (props) => ({
            ...props,
            color: 'gray.400',
            textTransform: 'uppercase',
            fontWeight: 'regular',
          }),
        }}
        {...props}
      />
    );
  },
);

export type { SelectInstance };

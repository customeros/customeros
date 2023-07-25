import { forwardRef, useCallback, useMemo } from 'react';
import {
  Select as _Select,
  Props,
  ControlProps,
  SelectInstance,
  chakraComponents,
} from 'chakra-react-select';

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

    const components = useMemo(
      () => ({
        Control,
        DropdownIndicator: () => null,
        ClearIndicator: () => null,

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
          container: (props) => ({ ...props, w: '100%', overflow: 'visible' }),
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

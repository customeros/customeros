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
          container: (props) => ({ ...props, w: '100%' }),
        }}
        {...props}
      />
    );
  },
);

export type { SelectInstance };

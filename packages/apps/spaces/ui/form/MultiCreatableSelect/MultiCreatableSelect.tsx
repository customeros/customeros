import React, { forwardRef, useCallback, useMemo } from 'react';

import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { OptionProps } from 'chakra-react-select';
import {
  AsyncCreatableSelect,
  chakraComponents,
  AsyncCreatableProps,
  ControlProps,
  MultiValueGenericProps,
} from '@ui/form/SyncSelect';
import { Tooltip } from '@ui/presentation/Tooltip';
import { multiCreatableSelectStyles } from '@ui/form/MultiCreatableSelect/styles';

interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  formId: string;
  customStyles?: any;
  withTooltip?: boolean;
  Option?: any;
}

export const MultiCreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ chakraStyles, ...props }, ref) => {
    const Control = useCallback(({ children, ...rest }: ControlProps) => {
      return (
        <chakraComponents.Control {...rest}>
          {children}
        </chakraComponents.Control>
      );
    }, []);
    const MultiValueLabel = useCallback((rest: MultiValueGenericProps<any>) => {
      if (props?.withTooltip) {
        return (
          <chakraComponents.MultiValueLabel {...rest}>
            <Tooltip
              label={rest.data.label.length > 0 ? rest.data.value : ''}
              placement='top'
            >
              {rest.data.label || rest.data.value}
            </Tooltip>
          </chakraComponents.MultiValueLabel>
        );
      }
      return (
        <chakraComponents.MultiValueLabel {...rest}>
          {rest.data.label || rest.data.value}
        </chakraComponents.MultiValueLabel>
      );
    }, []);

    const Option = useCallback(
      (rest: OptionProps<{ label: string; value: string }>) => {
        return (
          <chakraComponents.Option {...rest}>
            {rest.data.label || rest.data.value}
          </chakraComponents.Option>
        );
      },
      [],
    );

    const components = useMemo(
      () => ({
        Control,
        MultiValueLabel,
        Option: props?.Option || Option,
        DropdownIndicator: () => null,
        ClearIndicator: () => null,
      }),
      [Control, MultiValueLabel],
    );

    return (
      <AsyncCreatableSelect
        loadOptions={props?.loadOptions}
        variant='unstyled'
        focusBorderColor='transparent'
        ref={ref}
        components={components}
        tabSelectsValue={false}
        isMulti
        tagVariant='ghost'
        chakraStyles={
          props?.customStyles?.(chakraStyles) ||
          multiCreatableSelectStyles(chakraStyles)
        }
        {...props}
      />
    );
  },
);

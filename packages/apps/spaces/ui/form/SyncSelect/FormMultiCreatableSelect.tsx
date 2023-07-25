import React, { forwardRef, useCallback, useMemo, useState } from 'react';

import { SelectInstance, SelectProps } from '@ui/form/SyncSelect/Select';
import { useField } from 'react-inverted-form';
import {
  AsyncCreatableSelect,
  chakraComponents,
  AsyncCreatableProps,
  ControlProps,
  MultiValueGenericProps,
} from 'chakra-react-select';
import { Tooltip } from '@chakra-ui/react';

interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  formId: string;
}

const MultiCreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ ...props }, ref) => {
    const Control = useCallback(({ children, ...rest }: ControlProps) => {
      return (
        <chakraComponents.Control {...rest}>
          {children}
        </chakraComponents.Control>
      );
    }, []);
    const MultiValueLabel = (props: MultiValueGenericProps<any>) => {
      return (
        <chakraComponents.MultiValueLabel {...props}>
          <Tooltip label={props.data.value} placement='top'>
            {props.data.label}
          </Tooltip>
        </chakraComponents.MultiValueLabel>
      );
    };

    const components = useMemo(
      () => ({
        Control,
        MultiValueLabel,
        DropdownIndicator: () => null,
        ClearIndicator: () => null,
      }),
      [Control, MultiValueLabel],
    );

    return (
      <AsyncCreatableSelect
        loadOptions={props.loadOptions}
        variant='unstyled'
        focusBorderColor='transparent'
        ref={ref}
        components={components}
        tabSelectsValue={false}
        isMulti
        tagVariant='ghost'
        chakraStyles={{
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
          }),
          clearIndicator: (base) => ({
            ...base,
            background: 'transparent',
            color: 'transparent',
            display: 'none',
          }),

          multiValueRemove: (styles, { data }) => ({
            ...styles,
            visibility: 'hidden',
            width: '0px',
          }),

          container: (props) => ({
            ...props,
            minWidth: '300px',
            zIndex: 9999,
            overflow: 'visible',
            _focusVisible: { border: 'none !important' },
            _focus: { border: 'none !important' },
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
          input: (props) => ({
            ...props,
            color: 'gray.500',
            fontWeight: 'regular',
          }),
        }}
        {...props}
      />
    );
  },
);

export const FormMultiCreatableSelect = forwardRef<
  SelectInstance,
  FormSelectProps
>(({ name, formId, ...rest }, ref) => {
  const { getInputProps } = useField(name, formId);
  const { id, onChange, onBlur, value } = getInputProps();
  return (
    <MultiCreatableSelect
      ref={ref}
      id={id}
      formId={formId}
      name={name}
      value={value}
      onBlur={() => onBlur(value)}
      onChange={onChange}
      {...rest}
    />
  );
});

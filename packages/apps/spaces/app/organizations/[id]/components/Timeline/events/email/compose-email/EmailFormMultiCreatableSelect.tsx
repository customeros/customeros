import React, { forwardRef, useCallback, useMemo } from 'react';

import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { useField } from 'react-inverted-form';
import { GroupBase, ChakraStylesConfig } from 'chakra-react-select';
import {
  AsyncCreatableSelect,
  chakraComponents,
  AsyncCreatableProps,
  ControlProps,
  MultiValueGenericProps,
} from '@ui/form/SyncSelect';
import { Tooltip } from '@ui/presentation/Tooltip';
import { emailRegex } from '@organization/components/Timeline/events/email/utils';
import omit from 'lodash/omit';

interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  formId: string;
}

const MultiCreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ chakraStyles, ...props }, ref) => {
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
          <Tooltip
            label={props.data.label.length > 0 ? props.data.value : ''}
            placement='top'
          >
            {props.data.label || props.data.value}
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
            cursor: 'default',
          }),
          clearIndicator: (base) => ({
            ...base,
            background: 'transparent',
            color: 'transparent',
            display: 'none',
          }),

          multiValueRemove: (styles, { data }) => ({
            ...styles,
            // visibility: 'hidden',
          }),

          container: (props) => ({
            ...props,
            minWidth: '300px',
            width: '100%',
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
            maxHeight: '12rem',
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
        }}
        {...props}
      />
    );
  },
);

export const EmailFormMultiCreatableSelect = forwardRef<
  SelectInstance,
  FormSelectProps
>(({ name, formId, ...rest }, ref) => {
  const { getInputProps } = useField(name, formId);
  const { id, onChange, onBlur, value } = getInputProps();
  const handleBlur = (stringVal: string) => {
    if (stringVal && emailRegex.test(stringVal)) {
      onBlur([...value, { label: stringVal, value: stringVal }]);
      return;
    }
    onBlur(value);
  };

  return (
    <MultiCreatableSelect
      ref={ref}
      id={id}
      formId={formId}
      name={name}
      value={value}
      onBlur={(e) => handleBlur(e.target.value)}
      onChange={onChange}
      {...rest}
    />
  );
});

import { useField } from 'react-inverted-form';
import React, {
  useId,
  useRef,
  ChangeEvent,
  cloneElement,
  PropsWithChildren,
} from 'react';

import {
  chakra,
  CheckboxIcon,
  CheckboxProps,
  SystemStyleObject,
  useMultiStyleConfig,
} from '@chakra-ui/react';

export interface FormCheckboxProps extends Omit<CheckboxProps, 'onChange'> {
  name: string;
  formId: string;
  onChange?: (value: boolean) => void;
}

const rootStyles: SystemStyleObject = {
  cursor: 'pointer',
  display: 'inline-flex',
  alignItems: 'center',
  verticalAlign: 'top',
  position: 'relative',
};

const controlStyles: SystemStyleObject = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  verticalAlign: 'top',
  userSelect: 'none',
  flexShrink: 0,
};

export const FormCheckbox = (props: PropsWithChildren<FormCheckboxProps>) => {
  const {
    id,
    value,
    children,
    iconSize,
    iconColor,
    isIndeterminate,
    spacing = '0.5rem',
    fontSize = 'md',
    icon = <CheckboxIcon />,
    formId,
    name,
  } = props;
  const { getInputProps } = useField(name, formId);
  const { value: isChecked, onChange } = getInputProps();
  const _id = useId();
  const inputRef = useRef<HTMLInputElement>(null);
  const styles = useMultiStyleConfig('Checkbox', props);
  const clonedIcon = cloneElement(icon, {
    __css: {
      fontSize: iconSize,
      color: iconColor,
      ...styles.icon,
    },
    isIndeterminate: isIndeterminate,
    isChecked: isChecked,
  });

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    onChange(e.target?.checked);
  };

  return (
    <chakra.label
      htmlFor={id ?? _id}
      __css={{ ...rootStyles, ...styles.container }}
    >
      <chakra.input
        hidden
        value={value}
        ref={inputRef}
        id={id ?? _id}
        type='checkbox'
        checked={isChecked}
        onChange={handleChange}
      />

      <chakra.span
        __css={{
          ...controlStyles,
          ...styles.control,
        }}
      >
        {clonedIcon}
      </chakra.span>
      <chakra.span
        __css={{
          ...styles.label,
          marginStart: spacing,
          fontSize,
          lineHeight: 1,
        }}
      >
        {children}
      </chakra.span>
    </chakra.label>
  );
};

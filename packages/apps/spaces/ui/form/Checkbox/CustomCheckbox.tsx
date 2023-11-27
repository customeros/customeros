import { useId, useRef, ChangeEvent, cloneElement } from 'react';

import {
  chakra,
  CheckboxIcon,
  CheckboxProps,
  SystemStyleObject,
  useMultiStyleConfig,
} from '@chakra-ui/react';

export interface CustomCheckboxProps extends Omit<CheckboxProps, 'onChange'> {
  onChange?: (value: string) => void;
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

export const CustomCheckbox = (props: CustomCheckboxProps) => {
  const _id = useId();
  const inputRef = useRef<HTMLInputElement>(null);
  const styles = useMultiStyleConfig('Checkbox', props);
  const {
    id,
    value,
    children,
    onChange,
    isChecked,
    iconSize,
    iconColor,
    isIndeterminate,
    spacing = '0.5rem',
    icon = <CheckboxIcon />,
  } = props;

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    onChange?.(e.target.value);
  };

  const clonedIcon = cloneElement(icon, {
    __css: {
      fontSize: iconSize,
      color: iconColor,
      ...styles.icon,
    },
    isIndeterminate: isIndeterminate,
    isChecked: isChecked,
  });

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
      <chakra.span __css={{ marginStart: spacing, ...styles.label }}>
        {children}
      </chakra.span>
    </chakra.label>
  );
};

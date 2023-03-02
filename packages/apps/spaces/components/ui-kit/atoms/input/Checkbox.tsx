import React from 'react';
import {
  Checkbox as PrimereactCheckbox,
  CheckboxProps,
} from 'primereact/checkbox';

export const Checkbox: React.FC<CheckboxProps> = ({ checked, ...props }) => {
  return <PrimereactCheckbox {...props} checked={checked}></PrimereactCheckbox>;
};

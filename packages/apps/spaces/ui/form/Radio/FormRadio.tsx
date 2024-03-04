'use client';

import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  FormLabel,
  RadioProps,
  FormControl,
  VisuallyHidden,
  FormLabelProps,
} from '@chakra-ui/react';

import { Text } from '@ui/typography/Text';

import { Radio, RadioGroup } from './Radio';

export interface FromRadioProps extends RadioProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
}

//todo add visually hidden label - accessibility

export const FromRadioGroup = forwardRef(
  (
    {
      name,
      formId,
      label,
      isLabelVisible,
      labelProps,
      defaultValue,
      children,
      ...props
    }: FromRadioProps,
    ref,
  ) => {
    const { getInputProps, renderError, state } = useField(name, formId);

    return (
      <FormControl>
        {isLabelVisible ? (
          <FormLabel {...labelProps}>{label}</FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}

        <RadioGroup
          ref={ref}
          {...getInputProps()}
          {...props}
          isInvalid={state.meta?.meta?.hasError}
          autoComplete='off'
          data-1p-ignore
        >
          {children}
        </RadioGroup>
        {renderError((error) => (
          <Text fontSize='xs' color='error.500'>
            {error}
          </Text>
        ))}
      </FormControl>
    );
  },
);

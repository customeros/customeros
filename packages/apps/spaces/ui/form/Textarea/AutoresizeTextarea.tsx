import React, { forwardRef } from 'react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import {
  Textarea,
  FormLabel,
  FormControl,
  TextareaProps,
  VisuallyHidden,
} from '@chakra-ui/react';

import { InputGroup, InputLeftElement } from '../InputGroup';

export interface AutoresizeTextareaProps
  extends TextareaProps,
    Pick<
      TextareaAutosizeProps,
      'maxRows' | 'minRows' | 'onHeightChange' | 'cacheMeasurements'
    > {
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
}

export const AutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  AutoresizeTextareaProps
>(({ leftElement, isLabelVisible, label = '', ...props }, ref) => {
  return (
    <FormControl>
      {isLabelVisible ? (
        <FormLabel
          fontWeight={600}
          color={props?.color}
          fontSize='sm'
          mb={0}
          mt={2}
        >
          {label}
        </FormLabel>
      ) : (
        <VisuallyHidden>
          <FormLabel>{label}</FormLabel>
        </VisuallyHidden>
      )}
      <InputGroup>
        {leftElement && (
          <InputLeftElement w='4'>{leftElement}</InputLeftElement>
        )}
        <Textarea
          w='100%'
          ref={ref}
          minRows={1}
          minH='unset'
          resize='none'
          aria-label={label}
          overflow='hidden'
          as={ResizeTextarea}
          borderColor='transparent'
          color='gray.700'
          _hover={{
            borderColor: 'gray.300',
          }}
          _focusVisible={{
            borderColor: 'primary.500',
            boxShadow: 'unset',
          }}
          {...props}
        />
      </InputGroup>
    </FormControl>
  );
});

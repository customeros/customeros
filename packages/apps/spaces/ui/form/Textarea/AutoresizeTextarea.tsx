import { forwardRef } from 'react';

import { Textarea, TextareaProps } from '@chakra-ui/react';
import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { InputGroup, InputLeftElement } from '../InputGroup';

export interface AutoresizeTextareaProps
  extends TextareaProps,
    Pick<
      TextareaAutosizeProps,
      'maxRows' | 'minRows' | 'onHeightChange' | 'cacheMeasurements'
    > {
  leftElement?: React.ReactNode;
}

export const AutoresizeTextarea = forwardRef<
  HTMLTextAreaElement,
  AutoresizeTextareaProps
>(({ leftElement, ...props }, ref) => {
  return (
    <InputGroup>
      {leftElement && <InputLeftElement>{leftElement}</InputLeftElement>}
      <Textarea
        w='100%'
        ref={ref}
        minRows={1}
        minH='unset'
        resize='none'
        overflow='hidden'
        as={ResizeTextarea}
        _focusVisible={{
          borderColor: 'teal.500',
          boxShadow: 'unset',
        }}
        {...props}
      />
    </InputGroup>
  );
});

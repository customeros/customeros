import { useField } from 'react-inverted-form';
import { useRef, useState, useEffect } from 'react';

import {
  Flex,
  Editable,
  EditableProps,
  EditableInput,
  EditablePreview,
} from '@chakra-ui/react';

interface FormEditableInputProps extends EditableProps {
  name: string;
  formId: string;
}

export const FormEditableInput = ({
  name,
  formId,
  ...props
}: FormEditableInputProps) => {
  const [width, setWidth] = useState('0px');
  const spanRef = useRef<HTMLSpanElement>(null);
  const { getInputProps } = useField(name, formId);
  const inputProps = getInputProps();

  useEffect(() => {
    const spanWidth = spanRef.current?.offsetWidth ?? 0;
    setWidth(`${spanWidth}px`);
  }, [inputProps.value]);

  return (
    <>
      <Editable value={inputProps.value} width={width} {...props}>
        <EditablePreview />
        <EditableInput autoComplete='off' spellCheck='false' {...inputProps} />
        <Flex
          as='span'
          ref={spanRef}
          sx={{
            zIndex: -1,
            position: 'absolute',
            paddingRight: '0.5rem',
            height: '0px',
            display: 'inline-block',
            visibility: 'hidden',
            whiteSpace: 'pre',
          }}
        >
          {inputProps.value}
        </Flex>
      </Editable>
    </>
  );
};

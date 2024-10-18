import { IMaskMixinProps } from 'react-imask';
import React, { useRef, useState, useEffect } from 'react';

import { MaskElement } from 'imask';

import { InputProps } from '@ui/form/Input/Input';
import { MaskedInput } from '@ui/form/Input/MaskedInput';
type MaskedInputProps = IMaskMixinProps<MaskElement> & InputProps;

export const MaskedResizableInput = ({ ...props }: MaskedInputProps) => {
  const spanRef = useRef<HTMLSpanElement>(null);
  const [width, setWidth] = useState('10px');

  useEffect(() => {
    const measureWidth = () => {
      if (spanRef.current) {
        const spanWidth = spanRef.current?.offsetWidth ?? 10;

        setWidth(`${Math.max(spanWidth + 2, 2)}px`); // Add some padding and set a minimum width
      }
    };

    measureWidth();
  }, [props.value, props.defaultValue]);

  const handleAccept = (unmaskedValue: string) => {
    if (props.onChange) {
      props.onChange({
        target: { value: unmaskedValue },
      } as React.ChangeEvent<HTMLInputElement>);
    }
  };

  return (
    <>
      <span
        ref={spanRef}
        className={`z-[-1] absolute h-0 inline-block invisible`}
      >
        {props.value || props.defaultValue || props.placeholder || ''}
      </span>

      <MaskedInput
        size='xs'
        variant='unstyled'
        onAccept={handleAccept}
        style={{ ...props.style, width: width }}
        {...props}
      />
    </>
  );
};

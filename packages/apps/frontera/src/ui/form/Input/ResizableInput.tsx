import { useRef, useState, useEffect, forwardRef } from 'react';

import { Input, InputProps } from '@ui/form/Input/Input';

export const ResizableInput = forwardRef<HTMLInputElement, InputProps>(
  (props: InputProps, ref) => {
    const spanRef = useRef<HTMLSpanElement>(null);
    const [width, setWidth] = useState('10px');

    useEffect(() => {
      const measureWidth = () => {
        if (spanRef.current) {
          const spanWidth = spanRef.current?.offsetWidth ?? 0;

          setWidth(`${spanWidth}px`);
        }
      };

      measureWidth();
    }, [props.value, props.defaultValue]);

    return (
      <>
        <span
          ref={spanRef}
          className={`z-[-1] absolute h-0 inline-block invisible`}
        >
          {props.value || props.defaultValue || props.placeholder || ''}
        </span>

        <Input ref={ref} data-1p-ignore {...props} style={{ width: width }} />
      </>
    );
  },
);

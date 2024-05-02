import { useRef, useState, forwardRef, useLayoutEffect } from 'react';

import { Input, InputProps } from '@ui/form/Input/Input';

export const ResizableInput = forwardRef<HTMLInputElement, InputProps>(
  (props: InputProps, ref) => {
    const spanRef = useRef<HTMLSpanElement>(null);
    const [width, setWidth] = useState('10px');
    useLayoutEffect(() => {
      const measureWidth = () => {
        if (spanRef.current) {
          const spanWidth = spanRef.current?.offsetWidth ?? 10;
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
          {props.value || props.defaultValue || '0'}
        </span>

        <Input
          ref={ref}
          data-1p-ignore
          {...props}
          style={{ width: width, minWidth: width }}
        />
      </>
    );
  },
);

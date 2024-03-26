import { useRef, useState, useEffect, forwardRef } from 'react';

import { Input, InputProps } from '@ui/form/Input/Input2';

export const ResizableInput = forwardRef<HTMLInputElement, InputProps>(
  (props: InputProps, ref) => {
    const spanRef = useRef<HTMLSpanElement>(null);
    const [width, setWidth] = useState('0px');
    useEffect(() => {
      const measureWidth = () => {
        if (spanRef.current) {
          const spanWidth = spanRef.current?.offsetWidth ?? 0;
          setWidth(`${spanWidth}px`);
        }
      };
      measureWidth();
    }, [props.value]);

    return (
      <>
        <span
          ref={spanRef}
          className={`z-[-1] absolute h-0 inline-block invisible`}
        >
          {props.value}
        </span>

        <Input ref={ref} data-1p-ignore {...props} style={{ width: width }} />
      </>
    );
  },
);

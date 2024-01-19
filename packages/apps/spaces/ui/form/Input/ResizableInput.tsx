import { useRef, useState, useEffect, forwardRef } from 'react';

import { Flex, Input, InputProps } from '@chakra-ui/react';

export const ResizableInput = forwardRef<HTMLInputElement, InputProps>(
  (props: InputProps, ref) => {
    const spanRef = useRef<HTMLSpanElement>(null);
    const [width, setWidth] = useState('0px');

    useEffect(() => {
      const spanWidth = spanRef.current?.offsetWidth ?? 0;
      setWidth(`${spanWidth}px`);
    }, [props.value]);

    return (
      <>
        <Flex
          as='span'
          ref={spanRef}
          fontSize={props.fontSize ?? props.size}
          fontWeight={props.fontWeight}
          sx={{
            zIndex: -1,
            position: 'absolute',
            height: '0px',
            display: 'inline-block',
            visibility: 'hidden',
            whiteSpace: 'pre',
          }}
        >
          {props.value}
        </Flex>
        <Input ref={ref} w={width} {...props} />
      </>
    );
  },
);

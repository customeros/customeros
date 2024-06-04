import { useRef, useEffect } from 'react';

import { Input, InputProps } from '@ui/form/Input';

export const OrganizationNameInput = ({
  orgNameReadOnly,
  isLoading,
  ...rest
}: {
  rest: InputProps;
  isLoading: boolean;
  orgNameReadOnly: boolean;
}) => {
  const nameRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (nameRef.current?.value === 'Unnamed' && !isLoading) {
      nameRef.current?.focus();
      nameRef.current?.setSelectionRange(0, 7);
    }
  }, [nameRef, isLoading]);

  return (
    <>
      <Input
        className='font-semibold text-lg border-none overflow-hidden overflow-ellipsis'
        name='name'
        ref={nameRef}
        autoComplete='off'
        variant='unstyled'
        placeholder='Company name'
        disabled={orgNameReadOnly}
        size='xs'
        {...rest}
      />
    </>
  );
};

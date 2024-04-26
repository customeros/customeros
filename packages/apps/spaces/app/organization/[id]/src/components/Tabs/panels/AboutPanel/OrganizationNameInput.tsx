'use client';
import { useRef, useEffect } from 'react';

import { FormInput } from '@ui/form/Input/FormInput';

export const OrganizationNameInput = ({
  orgNameReadOnly,
  isLoading,
}: {
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
      <FormInput
        className='font-semibold text-lg border-none overflow-hidden overflow-ellipsis'
        name='name'
        ref={nameRef}
        autoComplete='off'
        variant='unstyled'
        placeholder='Company name'
        formId='organization-about'
        disabled={orgNameReadOnly}
        size='xs'
      />
    </>
  );
};

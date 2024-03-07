'use client';
import { useRef, useEffect } from 'react';

import { FormInput } from '@ui/form/Input';

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
        display='block'
        name='name'
        fontSize='lg'
        ref={nameRef}
        autoComplete='off'
        fontWeight='semibold'
        variant='unstyled'
        borderRadius='unset'
        placeholder='Company name'
        formId='organization-about'
        isReadOnly={orgNameReadOnly}
        textOverflow='ellipsis'
        overflow='hidden'
      />
    </>
  );
};

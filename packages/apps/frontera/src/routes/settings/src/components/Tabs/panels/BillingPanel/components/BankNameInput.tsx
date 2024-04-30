'use client';
import React, { useRef, useEffect } from 'react';

import { Metadata } from '@graphql/types';
import { FormInput } from '@ui/form/Input/FormInput';

export const BankNameInput = ({
  formId,
  metadata,
}: {
  formId: string;
  metadata: Metadata;
}) => {
  const nameRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    const wasEdited = metadata.created !== metadata.lastUpdated;
    const isDefaultName = /^[A-Z]{3}\saccount$/.test(
      nameRef.current?.value ?? '',
    );
    if (isDefaultName && !wasEdited) {
      nameRef.current?.focus();
      nameRef.current?.setSelectionRange(0, 11);
    }
  }, [nameRef, metadata.created, metadata.lastUpdated]);

  return (
    <>
      <FormInput
        ref={nameRef}
        className='text-md font-semibold'
        autoComplete='off'
        label='Bank Name'
        labelProps={{ className: 'hidden' }}
        placeholder='Bank name'
        name='bankName'
        formId={formId}
        variant={'unstyled'}
      />
    </>
  );
};

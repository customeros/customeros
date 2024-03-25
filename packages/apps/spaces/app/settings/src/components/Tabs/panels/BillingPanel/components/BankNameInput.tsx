'use client';
import React, { useRef, useEffect } from 'react';

import { Metadata } from '@graphql/types';
import { FormInput } from '@ui/form/Input';

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
        fontSize='md'
        fontWeight='semibold'
        autoComplete='off'
        label='Bank Name'
        placeholder='Bank name'
        name='bankName'
        formId={formId}
        border='none'
        _hover={{
          border: 'none',
        }}
        _focus={{
          border: 'none',
        }}
        _focusVisible={{
          border: 'none',
        }}
      />
    </>
  );
};

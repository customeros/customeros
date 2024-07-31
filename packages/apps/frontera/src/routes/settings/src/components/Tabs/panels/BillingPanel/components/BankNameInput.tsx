import { useRef, useEffect } from 'react';

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
        name='bankName'
        formId={formId}
        label='Bank Name'
        autoComplete='off'
        variant={'unstyled'}
        placeholder='Bank name'
        className='text-md font-semibold'
        labelProps={{ className: 'hidden' }}
      />
    </>
  );
};

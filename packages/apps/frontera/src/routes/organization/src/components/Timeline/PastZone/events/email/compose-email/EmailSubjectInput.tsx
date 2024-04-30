'use client';
import React, { FC } from 'react';

import { FormInput } from '@ui/form/Input/FormInput';

interface EmailSubjectInput {
  mt?: number;
  formId: string;
  fieldName: string;
}

export const EmailSubjectInput: FC<EmailSubjectInput> = ({
  fieldName,
  formId,
  mt = 0,
}) => {
  return (
    <div className='flex items-center flex-1' style={{ marginTop: mt }}>
      <span className='text-gray-700 font-semibold mr-1'>Subject:</span>
      <FormInput
        name={fieldName}
        formId={formId}
        className='text-gray-500 height-[5px] text-md'
        variant='unstyled'
        size='xs'
      />
    </div>
  );
};

import { FC } from 'react';

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
    <div style={{ marginTop: mt }} className='flex items-center flex-1'>
      <span className='text-gray-700 font-semibold mr-1'>Subject:</span>
      <FormInput
        size='xs'
        formId={formId}
        name={fieldName}
        variant='unstyled'
        className='text-gray-500 height-[5px] text-md'
      />
    </div>
  );
};

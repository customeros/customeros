import { forwardRef } from 'react';
import { SelectInstance } from 'react-select';

import { EmailFormMultiCreatableSelect } from '@shared/components/EmailMultiCreatableSelect';

interface EmailParticipantSelect {
  formId: string;
  entryType: string;
  fieldName: string;
  autofocus?: boolean;
  placeholder?: string;
}

export const EmailSelect = forwardRef<SelectInstance, EmailParticipantSelect>(
  (
    {
      entryType,
      fieldName,
      formId,
      autofocus = false,
      placeholder = 'Enter email',
    },
    ref,
  ) => {
    return (
      <div className='text-sm'>
        <label className='font-semibold text-sm'>{entryType}</label>
        <EmailFormMultiCreatableSelect
          ref={ref}
          name={fieldName}
          formId={formId}
          autoFocus={autofocus}
          placeholder={placeholder}
          navigateAfterAddingToPeople={true}
          noOptionsMessage={() => null}
          allowCreateWhileLoading={false}
          formatCreateLabel={(input) => {
            return input;
          }}
          getOptionLabel={(d) => {
            if (d?.__isNew__) {
              return `${d.label}`;
            }

            return `${d.label} - ${d.value}`;
          }}
        />
      </div>
    );
  },
);

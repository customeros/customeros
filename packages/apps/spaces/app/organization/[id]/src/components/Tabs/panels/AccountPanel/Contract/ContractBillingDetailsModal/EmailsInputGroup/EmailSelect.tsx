'use client';
import React, { FC } from 'react';

import { EmailFormMultiCreatableSelect } from '@shared/components/EmailMultiCreatableSelect';
import { emailRegex } from '@organization/src/components/Timeline/PastZone/events/email/utils';

interface EmailParticipantSelect {
  formId: string;
  entryType: string;
  fieldName: string;
  autofocus: boolean;
  placeholder?: string;
}

export const EmailSelect: FC<EmailParticipantSelect> = ({
  entryType,
  fieldName,
  formId,
  autofocus = false,
  placeholder = 'Enter email',
}) => {
  return (
    <div>
      <label className='font-semibold text-sm'>{entryType}</label>
      <EmailFormMultiCreatableSelect
        autoFocus={autofocus}
        name={fieldName}
        formId={formId}
        placeholder={placeholder}
        navigateAfterAddingToPeople={true}
        noOptionsMessage={() => null}
        allowCreateWhileLoading={false}
        formatCreateLabel={(input) => {
          return input;
        }}
        isValidNewOption={(input) => emailRegex.test(input)}
        getOptionLabel={(d) => {
          if (d?.__isNew__) {
            return `${d.label}`;
          }

          return `${d.label} - ${d.value}`;
        }}
      />
    </div>
  );
};

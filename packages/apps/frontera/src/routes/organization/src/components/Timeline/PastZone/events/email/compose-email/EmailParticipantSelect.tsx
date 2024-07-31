import { FC } from 'react';

import { emailRegex } from '@organization/components/Timeline/PastZone/events/email/utils';
import { EmailFormMultiCreatableSelect } from '@shared/components/EmailMultiCreatableSelect';

interface EmailParticipantSelect {
  formId: string;
  entryType: string;
  fieldName: string;
  autofocus: boolean;
}

export const EmailParticipantSelect: FC<EmailParticipantSelect> = ({
  entryType,
  fieldName,
  formId,
  autofocus = false,
}) => {
  return (
    <div className='flex items-baseline mb-[-1px] mt-0 flex-1 overflow-visible'>
      <span className='text-gray-700 font-semibold mr-1'>{entryType}:</span>
      <EmailFormMultiCreatableSelect
        formId={formId}
        name={fieldName}
        autoFocus={autofocus}
        noOptionsMessage={() => null}
        allowCreateWhileLoading={false}
        navigateAfterAddingToPeople={true}
        placeholder='Enter name or email...'
        isValidNewOption={(input) => emailRegex.test(input)}
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
};

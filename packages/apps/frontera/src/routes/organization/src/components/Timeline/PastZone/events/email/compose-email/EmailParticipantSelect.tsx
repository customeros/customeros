import { FC } from 'react';

import { EmailParticipantSelectUsecase } from '@domain/usecases/email-composer/email-participant-select.usecase.ts';

import { EmailFormMultiCreatableSelect } from '@shared/components/EmailMultiCreatableSelect';

interface EmailParticipantSelect {
  entryType: string;
  autofocus: boolean;
  emailParticipantUseCase: EmailParticipantSelectUsecase;
}

export const EmailParticipantSelect: FC<EmailParticipantSelect> = ({
  entryType,
  autofocus = false,
  emailParticipantUseCase,
}) => {
  return (
    <div className='flex  mb-[-1px] mt-0 flex-1 overflow-visible min-h-[28px] items-center'>
      <span className='text-gray-700 font-semibold text-sm'>{entryType}:</span>
      <EmailFormMultiCreatableSelect
        name={entryType}
        autoFocus={autofocus}
        navigateAfterAddingToPeople={true}
        placeholder='Enter name or email...'
        emailParticipantUseCase={emailParticipantUseCase}
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

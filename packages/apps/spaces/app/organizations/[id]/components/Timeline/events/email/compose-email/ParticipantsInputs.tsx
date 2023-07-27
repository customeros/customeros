import React, { FC } from 'react';
import { EmailParticipantSelect } from '@organization/components/Timeline/events/email/compose-email/EmailParticipantSelect';

interface ParticipantInputsProps {
  showCC: boolean;
  showBCC: boolean;
}
export const ParticipantInputs: FC<ParticipantInputsProps> = ({ showCC, showBCC }) => (
  <>
    <EmailParticipantSelect
      formId='compose-email-preview'
      fieldName='to'
      entryType='To'
    />

    {showCC && (
      <EmailParticipantSelect
        formId='compose-email-preview'
        fieldName='cc'
        entryType='CC'
      />
    )}
    {showBCC && (
      <EmailParticipantSelect
        formId='compose-email-preview'
        fieldName='Bcc'
        entryType='BCC'
      />
    )}
  </>
);

import React from 'react';

import { EmailValidationDetails } from '@graphql/types';
import { SimpleValidationIndicator } from '@ui/presentation/validation/simple-validation-indicator';
import { VALIDATION_MESSAGES } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/utils.ts';
function isValidEmail(email: string) {
  // Regular expression for validating an email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  // Test the email against the regex
  return emailRegex.test(email);
}

interface EmailCellProps {
  email: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailCell: React.FC<EmailCellProps> = ({
  email,
  validationDetails,
}) => {
  if (!email) return <p className='text-gray-400'>Unknown</p>;
  const getMessages = () => {
    if (!validationDetails) return [];
    const { validated, isReachable, isValidSyntax } = validationDetails;
    if (!validated && !isValidEmail(email))
      return [VALIDATION_MESSAGES.isValidSyntax.message];

    if (validated && isValidSyntax === false) {
      return [VALIDATION_MESSAGES.isValidSyntax.message];
    }
    if (
      validated &&
      isReachable &&
      (VALIDATION_MESSAGES.isReachable.condition as Array<string>).includes(
        isReachable,
      )
    ) {
      return [VALIDATION_MESSAGES.isReachable.message];
    }

    return [];
  };

  return (
    <div className='flex align-middle'>
      <p className='max-w-[140px] overflow-ellipsis overflow-hidden'>{email}</p>
      <SimpleValidationIndicator
        errorMessages={getMessages()}
        showValidationMessage={true}
        isLoading={false}
      />
    </div>
  );
};

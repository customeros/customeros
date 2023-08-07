import React, { useEffect, useState } from 'react';
import { EmailValidationDetails } from '@spaces/graphql';
import { VALIDATION_MESSAGES } from '@spaces/molecules/communication-details/utils';
import { validateEmail } from '@spaces/molecules/communication-details/useEmailValidation';
import { useTenantName } from '@spaces/hooks/useTenant';
import { useRecoilValue } from 'recoil';
import { tenantName } from '../../../../state/userData';
import { SimpleValidationIndicator } from '@spaces/ui/presentation/validation/simple-validation-indicator';

interface Props {
  isEditMode: boolean;
  email: string;
  showValidationMessage: boolean;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailValidationMessage = ({
  isEditMode,
  showValidationMessage,
  email,
  validationDetails,
}: Props) => {
  const [isLoading, setIsLoading] = useState(!validationDetails);
  const [validationData, setValidationData] = useState<
    EmailValidationDetails | null | undefined
  >(validationDetails);

  useTenantName();
  const tenant = useRecoilValue(tenantName);

  useEffect(() => {
    if (!validationDetails) {
      validateEmail({ email, tenant }).then((result) => {
        setIsLoading(false);
        if (result) {
          setValidationData(result);
        }
      });
    }
  }, [email]);

  if (!validationData && !isLoading) {
    return null;
  }

  const getMessages = () => {
    const messages: Array<string> = [];
    if (!validationData) return messages;
    const { validated, isReachable, ...input } = validationData;

    for (const key in input) {
      if (
        //@ts-expect-error fixme
        input[key] !== null &&
        Object.prototype.hasOwnProperty.call(VALIDATION_MESSAGES, key) &&
        //@ts-expect-error fixme
        VALIDATION_MESSAGES[key]?.condition === input[key]
      ) {
        messages.push(VALIDATION_MESSAGES[key].message);
      }
    }
    if (
      isReachable &&
      (VALIDATION_MESSAGES.isReachable.condition as Array<string>).includes(
        isReachable,
      )
    ) {
      messages.push(VALIDATION_MESSAGES.isReachable.message);
    }
    return messages;
  };

  return (
    <SimpleValidationIndicator
      errorMessages={getMessages()}
      showValidationMessage={showValidationMessage}
      isLoading={isLoading}
    />
  );
};

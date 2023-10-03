import React, { useEffect, useState } from 'react';
import { EmailValidationDetails } from '@graphql/types';
import { SimpleValidationIndicator } from '@ui/presentation/validation/simple-validation-indicator';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  validateEmail,
  VALIDATION_MESSAGES,
} from '@organization/src/components/Tabs/panels/PeoplePanel/ContactCard/utils';

interface Props {
  email: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailValidationMessage = ({ email, validationDetails }: Props) => {
  const client = getGraphQLClient();
  const [isLoading, setIsLoading] = useState(!validationDetails);
  const [validationData, setValidationData] = useState<
    EmailValidationDetails | null | undefined
  >(validationDetails);

  const { data: tenantNameQuery } = useTenantNameQuery(client);

  useEffect(() => {
    if (!validationDetails && tenantNameQuery?.tenant) {
      validateEmail({ email, tenant: tenantNameQuery?.tenant }).then(
        (result) => {
          setIsLoading(false);
          if (result) {
            setValidationData(result);
          }
        },
      );
    }
  }, [email, tenantNameQuery?.tenant]);

  if (!validationData && !isLoading) {
    return null;
  }
  const getMessages = () => {
    const messages: Array<string> = [];
    if (!validationData) return messages;
    const { validated, isReachable, ...input } = validationData;

    for (const key in input) {
      if (
        //@ts-expect-error improve type
        input[key] !== null &&
        Object.prototype.hasOwnProperty.call(VALIDATION_MESSAGES, key) &&
        //@ts-expect-error improve type
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
      showValidationMessage={true}
      isLoading={isLoading}
    />
  );
};

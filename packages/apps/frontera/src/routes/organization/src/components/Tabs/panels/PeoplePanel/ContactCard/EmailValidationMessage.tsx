import { useState, useEffect } from 'react';

import { EmailValidationDetails } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';
import { SimpleValidationIndicator } from '@ui/presentation/validation/simple-validation-indicator';

import { validateEmail, VALIDATION_MESSAGES } from './utils';

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
    if (!validationData) return [];
    const { validated, isReachable, isValidSyntax } = validationData;

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
    <SimpleValidationIndicator
      errorMessages={getMessages()}
      showValidationMessage={true}
      isLoading={isLoading}
    />
  );
};

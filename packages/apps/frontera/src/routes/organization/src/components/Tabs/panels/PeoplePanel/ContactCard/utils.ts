import axios from 'axios';

import { EmailValidationDetails } from '@graphql/types';

interface Props {
  email: string;
  tenant: string;
}
export const validateEmail = ({ email, tenant }: Props) => {
  return axios
    .post(
      `/validation-api/validateEmail`,
      { email },
      {
        headers: {
          'X-Openline-TENANT': tenant,
          'Content-Type': 'application/json',
        },
      },
    )
    .then((response: { data: EmailValidationDetails | null | undefined }) => {
      return response?.data ?? null;
    })
    .catch((_) => {
      return;
    });
};

export interface MappedObject {
  [key: string]: {
    message: string;
    condition: boolean | Array<string>;
  };
}

export interface InputObject {
  [key: string]: boolean | null | string;
}

export const VALIDATION_MESSAGES: MappedObject = {
  isReachable: {
    condition: ['invalid'],
    message: 'This mailbox is not reachable.',
  },
  isValidSyntax: {
    condition: false,
    message: 'This email address appears to be invalid',
  },
};

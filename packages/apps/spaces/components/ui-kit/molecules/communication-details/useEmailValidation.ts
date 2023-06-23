import axios from 'axios';
import { EmailValidationDetails } from '@spaces/graphql';

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
    .catch((e) => {
      console.error('Validation error', e);
      return;
    });
};

import { ApolloError } from '@apollo/client';
import {
  useGetContactNameByEmailLazyQuery,
  useGetContactNameByIdLazyQuery,
  useGetContactNameByPhoneNumberLazyQuery,
} from './types';
import { useEffect, useState } from 'react';
import { getContactDisplayName } from '../../utils';

interface Props {
  phoneNumber?: string;
  email?: string;
  id?: string;
}

interface Result {
  contactName: string;
  loading: boolean;
  error: ApolloError | undefined;
}
export const useContactNameLazy = ({
  phoneNumber,
  email,
  id,
}: Props): Result => {
  const [getContactNameByEmail, { loading, error }] =
    useGetContactNameByEmailLazyQuery();
  const [getContactNameByPhoneNumber] =
    useGetContactNameByPhoneNumberLazyQuery();
  const [getContactNameById] = useGetContactNameByIdLazyQuery();
  const [contactName, setContactName] = useState('Unnamed');

  const handleGetContactNameByPhoneNumber = async (e164: string) => {
    const name = await getContactNameByPhoneNumber({ variables: { e164 } });
    if (
      name.data?.contact_ByPhone?.name ||
      name.data?.contact_ByPhone?.firstName ||
      name.data?.contact_ByPhone?.lastName
    ) {
      const displayName = getContactDisplayName(name.data?.contact_ByPhone);
      setContactName(displayName);
    }
  };
  const handleGetContactNameByEmail = async (email: string) => {
    const name = await getContactNameByEmail({ variables: { email } });
    if (
      name.data?.contact_ByEmail?.name ||
      name.data?.contact_ByEmail?.firstName ||
      name.data?.contact_ByEmail?.lastName
    ) {
      const displayName = getContactDisplayName(name.data?.contact_ByEmail);
      setContactName(displayName);
    }
  };
  const handleGetContactNameById = async (id: string) => {
    const name = await getContactNameById({ variables: { id } });
    if (
      name.data?.contact?.name ||
      name.data?.contact?.firstName ||
      name.data?.contact?.lastName
    ) {
      const displayName = getContactDisplayName(name.data?.contact);
      setContactName(displayName);
    }
  };

  useEffect(() => {
    if (phoneNumber) {
      handleGetContactNameByPhoneNumber(phoneNumber);
    }
    if (email) {
      handleGetContactNameByEmail(email);
    }
    if (id) {
      handleGetContactNameById(id);
    }
  }, [id, phoneNumber, email]);

  return {
    contactName,
    loading,
    error,
  };
};

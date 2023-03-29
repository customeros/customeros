import React, { useEffect, useState } from 'react';
import { CallParty } from '../../atoms';
import { useContactNameFromPhoneNumber } from '../../../../hooks/useContact';
import { getContactDisplayName } from '../../../../utils';
import { useUser } from '../../../../hooks/useUser';

export const ConversationPartyEmail = ({ email }: { email: string }) => {
  const { data } = useUser({
    email,
  });
  const [initials, setInitials] = useState<Array<string>>([]);
  const [name, setName] = useState(email);
  useEffect(() => {
    if (!initials.length && data) {
      const name = getContactDisplayName(data);
      const initials = (name !== 'Unnamed' ? name : '').split(' ');
      if (initials.length) {
        setInitials(initials);
      }
      setName(name);
    }
  }, [data?.id, email]);

  return <CallParty direction='right' name={name} />;
};
export const ConversationPartyPhone = ({ tel }: { tel: string }) => {
  const { data } = useContactNameFromPhoneNumber({
    phoneNumber: tel,
  });

  const [initials, setInitials] = useState<Array<string>>([]);
  const [name, setName] = useState(tel);

  useEffect(() => {
    if (!initials.length && data) {
      const name = getContactDisplayName(data);
      const initials = (name !== 'Unnamed' ? name : '').split(' ');
      if (initials.length) {
        setInitials(initials);
      }

      setName(name);
    }
  }, [data?.id, tel, initials.length]);

  return <CallParty direction='left' name={name} />;
};

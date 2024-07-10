import React from 'react';

import { parsePhoneNumber, isPossiblePhoneNumber } from 'libphonenumber-js';

interface PhoneCellProps {
  phone: string;
}

export const PhoneCell: React.FC<PhoneCellProps> = ({ phone }) => {
  if (!phone) return;

  if (!isPossiblePhoneNumber(phone)) return <p>{phone}</p>;
  const parsedPhoneNumber = parsePhoneNumber(phone);

  return (
    <div className='flex align-middle'>
      <p> {parsedPhoneNumber.formatNational()}</p>
    </div>
  );
};

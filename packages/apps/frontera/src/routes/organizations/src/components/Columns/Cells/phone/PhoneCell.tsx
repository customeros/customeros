import React from 'react';

import { parsePhoneNumber } from 'libphonenumber-js';

interface PhoneCellProps {
  phone: string;
}

export const PhoneCell: React.FC<PhoneCellProps> = ({ phone }) => {
  if (!phone) return <p className='text-gray-400'>Unknown</p>;

  const parsedPhoneNumber = parsePhoneNumber(phone);

  return (
    <div className='flex align-middle'>
      <p> {parsedPhoneNumber.formatNational()}</p>
    </div>
  );
};

import React from 'react';

import { observer } from 'mobx-react-lite';
import flags from '@assets/countries/flags.json';
import countries from '@assets/countries/countries.json';

interface ContactNameCellProps {
  countryCode: string;
}

export const CountryCell: React.FC<ContactNameCellProps> = observer(
  ({ countryCode }) => {
    if (!countryCode) {
      return <div className='text-gray-400'>Unknown</div>;
    }

    const flag = flags[countryCode?.toLowerCase() as keyof typeof flags];
    const country = countries.find(
      (d) => d.alpha2 === countryCode.toLowerCase(),
    )?.name;

    return (
      <div className='flex items-center'>
        <img
          src={flag}
          alt={country}
          className='rounded-full mr-2'
          style={{ clipPath: 'circle(35%)' }}
        />
        <span className='overflow-hidden overflow-ellipsis whitespace-nowrap'>
          {country}
        </span>
      </div>
    );
  },
);

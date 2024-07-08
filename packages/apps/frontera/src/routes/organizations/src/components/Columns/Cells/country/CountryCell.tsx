import React from 'react';

import { observer } from 'mobx-react-lite';
import countries from '@assets/countries/countries.json';

import { flags } from '@ui/media/flags';

interface ContactNameCellProps {
  countryCode: string;
}

export const CountryCell: React.FC<ContactNameCellProps> = observer(
  ({ countryCode }) => {
    if (!countryCode) {
      return <div className='text-gray-400'>Unknown</div>;
    }

    const country = countries.find(
      (d) => d.alpha2 === countryCode.toLowerCase(),
    )?.name;

    return (
      <div className='flex items-center'>
        {flags[countryCode]}
        <span className='ml-2 overflow-hidden overflow-ellipsis whitespace-nowrap'>
          {country}
        </span>
      </div>
    );
  },
);

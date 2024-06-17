import {
  SingleValueProps,
  SelectComponentsConfig,
  components as selectComponents,
} from 'react-select';

import { DateTimeUtils } from '@utils/date';
import { SelectProps } from '@ui/form/Select';
import { Clock } from '@ui/media/icons/Clock';
import { Select } from '@ui/form/Select/Select';

const SingleValue = (props: SingleValueProps) => {
  const rawTimezone = props.children as string;
  const timezone = rawTimezone?.includes('UTC')
    ? rawTimezone.split(' ')[0]
    : rawTimezone;

  const time = DateTimeUtils.convertToTimeZone(
    new Date(),
    DateTimeUtils.defaultTimeFormatString,
    timezone,
  );
  const value = `${time} local time`;

  return (
    <selectComponents.SingleValue {...props}>
      <span className='text-gray-700 line-clamp-1'>
        {value}
        {` `}
        <span className='text-gray-500'>• {timezone}</span>
      </span>
    </selectComponents.SingleValue>
  );
};

const components = {
  SingleValue,
};

export const TimezoneSelect = ({ ...props }: SelectProps) => {
  return (
    <Select
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      components={components as SelectComponentsConfig<any, any, any>}
      leftElement={<Clock className='text-gray-500 mr-3' />}
      {...props}
    />
  );
};

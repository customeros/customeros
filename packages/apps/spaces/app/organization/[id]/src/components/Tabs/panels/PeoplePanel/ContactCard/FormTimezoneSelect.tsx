import {
  SingleValueProps,
  SelectComponentsConfig,
  components as selectComponents,
} from 'react-select';

import { Icons } from '@ui/media/Icon';
import { SelectProps } from '@ui/form/Select';
import { DateTimeUtils } from '@spaces/utils/date';
import { FormSelect } from '@ui/form/Select/FormSelect';

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

interface FormTimezoneSelectProps extends SelectProps {
  name: string;
  formId: string;
}

export const FormTimezoneSelect = ({ ...props }: FormTimezoneSelectProps) => {
  return (
    <FormSelect
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      components={components as SelectComponentsConfig<any, any, any>}
      leftElement={<Icons.Clock color='gray.500' mr='3' />}
      {...props}
    />
  );
};

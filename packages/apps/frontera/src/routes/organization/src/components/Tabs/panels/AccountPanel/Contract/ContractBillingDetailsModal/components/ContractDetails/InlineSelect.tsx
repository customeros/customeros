import { FC } from 'react';

import {
  Select,
  SelectProps,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select';

interface InlineSelectProps extends SelectProps {
  id: string;
  name: string;
  label: string;
  placeholder: string;
}

export const InlineSelect: FC<InlineSelectProps> = ({
  label,
  name,
  placeholder,
  options,
  id,
  onChange,
  onBlur,
  value,
  ...rest
}) => {
  const formSelectClassNames =
    'text-base inline min-h-1 max-h-3 border-none hover:border-none focus:border-none w-fit ml-1 mt-0 underline text-gray-500 hover:text-gray-700 focus:text-gray-700 min-w-[max-content]';

  const selectedOption = options?.find((option) => option.value === value);

  return (
    <div className='w-full'>
      <label className='absolute top-[-999999px]'>{label}</label>

      <Select
        name={name}
        options={options}
        onChange={onChange}
        value={selectedOption}
        defaultValue={selectedOption}
        onBlur={() => onBlur?.(value)}
        className={formSelectClassNames}
        classNames={{
          ...rest.classNames,
          container: () =>
            getContainerClassNames(
              'text-gray-500 text-base hover:text-gray-700 focus:text-gray-700 min-w-fit w-max-content z-10',
              undefined,
              { size: 'xs' },
            ),
          menuList: () => getMenuListClassNames('min-w-[120px]'),
        }}
        {...rest}
      />
    </div>
  );
};

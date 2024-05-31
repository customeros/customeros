import { FC } from 'react';

import { SelectOption } from '@ui/utils/types';
import { FormSelect } from '@ui/form/Select/FormSelect';
import {
  SelectProps,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select';

interface InlineSelectProps extends SelectProps {
  name: string;
  label: string;
  formId: string;
  placeholder: string;
  options: Array<SelectOption<unknown>>;
}

export const InlineSelect: FC<InlineSelectProps> = ({
  formId,
  label,
  name,
  placeholder,
  options,
  ...rest
}) => {
  const formSelectClassNames =
    'text-base inline min-h-1 max-h-3 border-none hover:border-none focus:border-none w-fit ml-1 mt-0 underline text-gray-500 hover:text-gray-700 focus:text-gray-700 min-w-[max-content]';

  return (
    <FormSelect
      className={formSelectClassNames}
      label={label}
      placeholder={placeholder}
      name={name}
      formId={formId}
      defaultMenuIsOpen={true}
      options={options}
      size='xs'
      classNames={{
        ...rest.classNames,
        container: () =>
          getContainerClassNames(
            'text-gray-500 text-base hover:text-gray-700 focus:text-gray-700 min-w-fit w-max-content z-10',
            { size: 'xs' },
          ),
        menuList: () => getMenuListClassNames('min-w-[120px]'),
      }}
      {...rest}
    />
  );
};

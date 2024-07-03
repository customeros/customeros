import { Select, SelectProps } from '@ui/form/Select';

interface MultiSelectFilterProps extends SelectProps {
  label: string;
  description: string;
  icon: React.ReactNode;
}

export const MultiSelectFilter = ({
  icon,
  label,
  description,
  ...rest
}: MultiSelectFilterProps) => {
  return (
    <div className='flex items-center w-full'>
      <div className='flex items-center flex-1'>
        {icon}
        <p className='font-medium'>
          {label}
          <span className='font-normal ml-1'>{description}</span>
        </p>
      </div>
      <div className='flex-1'>
        <Select isMulti isClearable={false} {...rest} />
      </div>
    </div>
  );
};

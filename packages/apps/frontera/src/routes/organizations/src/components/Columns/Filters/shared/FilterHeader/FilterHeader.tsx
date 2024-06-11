import { useRef, useEffect } from 'react';

import { Switch } from '@ui/form/Switch/Switch';

interface FilterHeaderProps {
  isChecked: boolean;
  onToggle: () => void;
  onDisplayChange: () => void;
}

export const FilterHeader = ({
  onToggle,
  isChecked,
  onDisplayChange,
}: FilterHeaderProps) => {
  const timeout = useRef<NodeJS.Timeout>();

  const handleChange = () => {
    onDisplayChange();

    if (timeout.current) {
      clearTimeout(timeout.current);
    }

    timeout.current = setTimeout(() => {
      onToggle();
    }, 0);
  };

  useEffect(() => {
    return () => {
      timeout.current && clearTimeout(timeout.current);
    };
  }, []);

  return (
    <div className='mb-3 flex items-center justify-between'>
      <span className='text-sm font-medium'>Filter</span>
      <Switch
        size='sm'
        colorScheme='primary'
        isChecked={isChecked}
        onChange={handleChange}
      />
    </div>
  );
};

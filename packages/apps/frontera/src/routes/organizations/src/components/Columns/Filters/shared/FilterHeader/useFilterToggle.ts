import { useState } from 'react';

interface FilterToggleOptions {
  defaultValue?: boolean;
  onToggle?: (setIsActive: (value: boolean) => void) => void;
}

export const useFilterToggle = ({
  onToggle,
  defaultValue,
}: FilterToggleOptions) => {
  const [isActive, setIsActive] = useState<boolean>(
    () => defaultValue ?? false,
  );

  const handleClick = (value?: boolean) => {
    setIsActive((prev) => value ?? !prev);
  };

  const handleChange = () => {
    onToggle?.(setIsActive);
  };

  return {
    isActive,
    setIsActive,
    handleClick,
    handleChange,
  };
};

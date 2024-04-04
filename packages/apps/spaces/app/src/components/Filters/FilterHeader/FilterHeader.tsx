import { useRef, useEffect } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { Text } from '@ui/typography/Text';

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
    }, 250);
  };

  useEffect(() => {
    return () => {
      timeout.current && clearTimeout(timeout.current);
    };
  }, []);

  return (
    <Flex
      mb='3'
      flexDir='row'
      alignItems='center'
      justifyContent='space-between'
    >
      <Text fontSize='sm' fontWeight='medium'>
        Filter
      </Text>
      <Switch
        size='sm'
        colorScheme='primary'
        isChecked={isChecked}
        onChange={handleChange}
      />
    </Flex>
  );
};

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormResizableInput } from '@ui/form/Input';

interface MilestoneDurationProps {
  formId: string;
  defaultValue: number;
  isActiveItem?: boolean;
  isMilestoneOpen?: boolean;
  onToggleMilestone: () => void;
}

export const MilestoneDuration = ({
  formId,
  defaultValue,
  isActiveItem,
  isMilestoneOpen,
  onToggleMilestone,
}: MilestoneDurationProps) => {
  const inputId = `${formId}-duration-input`;

  const handleClick = () => {
    if (!isMilestoneOpen) {
      onToggleMilestone();
    }
  };

  return (
    <Flex align='center' gap='1'>
      <Text
        as='label'
        fontSize='sm'
        color='gray.500'
        htmlFor={inputId}
        whiteSpace='nowrap'
      >
        Max duration:
      </Text>
      {isActiveItem ? (
        <Text fontSize='sm'>{defaultValue}</Text>
      ) : (
        <FormResizableInput
          min={1}
          size='sm'
          id={inputId}
          type='number'
          name='duration'
          formId={formId}
          variant='unstyled'
          borderRadius='unset'
          onClick={handleClick}
        />
      )}
      <Text fontSize='sm' color='gray.500' whiteSpace='nowrap'>
        {defaultValue === 1 ? 'day' : 'days'}
      </Text>
    </Flex>
  );
};

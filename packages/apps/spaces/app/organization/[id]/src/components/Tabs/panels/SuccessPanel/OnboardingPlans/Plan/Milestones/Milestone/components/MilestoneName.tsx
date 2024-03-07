import { useRef } from 'react';

import { FormInput } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';

interface MilestoneNameProps {
  formId: string;
  isLast?: boolean;
  defaultValue: string;
  isActiveItem?: boolean;
  isMilestoneOpen?: boolean;
  onToggleMilestone?: () => void;
}

export const MilestoneName = ({
  formId,
  defaultValue,
  isActiveItem,
  isMilestoneOpen,
  onToggleMilestone,
}: MilestoneNameProps) => {
  const nameInputRef = useRef<HTMLInputElement>(null);

  const handleClick = () => {
    if (!isMilestoneOpen) {
      onToggleMilestone?.();
    }
  };

  if (isActiveItem)
    return (
      <Text fontWeight='medium' w='full'>
        {defaultValue}
      </Text>
    );

  return (
    <FormInput
      name='name'
      formId={formId}
      ref={nameInputRef}
      variant='unstyled'
      fontWeight='medium'
      borderRadius='unset'
      onClick={handleClick}
      placeholder='Milestone name'
    />
  );
};

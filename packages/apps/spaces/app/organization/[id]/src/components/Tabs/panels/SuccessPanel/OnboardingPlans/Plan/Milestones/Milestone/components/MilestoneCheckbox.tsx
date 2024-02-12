import { ChangeEvent } from 'react';
import { useField } from 'react-inverted-form';

import { OnboardingPlanMilestoneStatus } from '@graphql/types';

import { StatusCheckbox } from '../StatusCheckbox';

interface MilestoneCheckboxProps {
  formId: string;
  readOnly?: boolean;
  colorScheme: string;
  showCustomIcon?: boolean;
  onToggleMilestone?: () => void;
}

export const MilestoneCheckbox = ({
  formId,
  readOnly,
  colorScheme,
  showCustomIcon,
  onToggleMilestone,
}: MilestoneCheckboxProps) => {
  const { getInputProps } = useField('statusDetails', formId);
  const { value, onChange, onBlur, ...inputProps } = getInputProps();

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (readOnly) {
      onToggleMilestone?.();

      return;
    }

    onChange?.({
      ...value,
      status: e.target.checked ? 'DONE' : 'NOT_STARTED',
      updatedAt: new Date().toISOString(),
    });
  };

  return (
    <StatusCheckbox
      mr='2'
      size='md'
      onChange={handleChange}
      colorScheme={colorScheme}
      showCustomIcon={showCustomIcon}
      isChecked={[
        OnboardingPlanMilestoneStatus.Done,
        OnboardingPlanMilestoneStatus.DoneLate,
      ].includes(value.status as unknown as OnboardingPlanMilestoneStatus)}
      {...inputProps}
    />
  );
};

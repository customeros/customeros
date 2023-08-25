import { useRef } from 'react';

import { FormSelect } from '@ui/form/SyncSelect';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';

interface FormRoleSelectProps {
  name: string;
  formId: string;
  isFocused: boolean;
  displayValue: string;
  placeholder?: string;
  isCardOpen?: boolean;
  setIsFocused: (isFocused: boolean) => void;
}

export const FormRoleSelect = ({
  name,
  formId,
  isFocused,
  isCardOpen,
  placeholder,
  displayValue,
  setIsFocused,
}: FormRoleSelectProps) => {
  const ref = useRef<HTMLDivElement>(null);

  useOutsideClick({
    ref,
    handler: () => setIsFocused(false),
  });

  if (isFocused) {
    return (
      <span onClick={(e) => e.stopPropagation()} ref={ref}>
        <FormSelect
          isMulti
          autoFocus
          menuIsOpen
          name={name}
          options={[
            { value: 'Decision Maker', label: 'Decision Maker' },
            { value: 'Influencer', label: 'Influencer' },
            { value: 'User', label: 'User' },
            { value: 'Stakeholder', label: 'Stakeholder' },
            { value: 'Gatekeeper', label: 'Gatekeeper' },
            { value: 'Champion', label: 'Champion' },
            { value: 'Interogator', label: 'Interogator' },
            { value: 'Data Owner', label: 'Data Owner' },
          ]}
          formId={formId}
          placeholder='Role'
        />
      </span>
    );
  }

  return (
    <Text
      cursor='text'
      color={displayValue ? 'gray.500' : 'gray.400'}
      onClick={(e) => {
        if (isCardOpen) {
          e.stopPropagation();
        }
        setIsFocused(true);
      }}
      borderBottom='1px solid transparent'
      transition='border-color 0.2s ease-in-out'
      _hover={{
        borderColor: 'gray.300',
      }}
    >
      {displayValue || placeholder}
    </Text>
  );
};

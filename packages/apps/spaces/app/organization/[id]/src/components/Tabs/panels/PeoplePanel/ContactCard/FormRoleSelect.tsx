import { useRef } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
import { FormSelect } from '@ui/form/SyncSelect';
import { SelectOption } from '@shared/types/SelectOptions';

import { RoleTag, getTagColorScheme } from './RoleTag';

interface FormRoleSelectProps {
  name: string;
  formId: string;
  isFocused: boolean;
  placeholder?: string;
  isCardOpen?: boolean;
  data: SelectOption<string>[];
  setIsFocused: (isFocused: boolean) => void;
}

export const FormRoleSelect = ({
  name,
  formId,
  isFocused,
  isCardOpen,
  placeholder,
  data,
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
            { value: 'Data Owner', label: 'Data Owner' },
          ]}
          formId={formId}
          placeholder='Role'
          chakraStyles={{
            multiValue: (props, data) => {
              const colorScheme = (() => getTagColorScheme(data.data.label))();

              return {
                ...props,
                fontSize: 'xs',
                fontWeight: 'normal',
                color: `${[colorScheme]}.700`,
                border: '1px solid',
                borderColor: `${[colorScheme]}.200`,
                backgroundColor: `${[colorScheme]}.50`,

                '& div[role="button"]': {
                  position: 'relative',
                  background: 'transparent',
                  outline: 'none',
                  marginInlineStart: '2px',
                  display: 'initial',
                  boxShadow: 'none !important',
                },
                '& div[data-focus="true"]': {
                  opacity: 1,
                },
              };
            },
          }}
        />
      </span>
    );
  }

  if (!data.length) {
    return (
      <Text
        cursor='text'
        color={'gray.400'}
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
        {placeholder}
      </Text>
    );
  }

  return (
    <Flex
      gap={1}
      mt={2}
      pb={2}
      flexWrap='wrap'
      onClick={(e) => {
        if (isCardOpen) {
          e.stopPropagation();
        }
        setIsFocused(true);
      }}
    >
      {data.map((e) => (
        <RoleTag key={e.label} label={e.label} />
      ))}
    </Flex>
  );
};

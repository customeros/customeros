import { useRef } from 'react';

import { colors } from '@ui/theme/colors';
import { useOutsideClick } from '@ui/utils';
import { FormSelect } from '@ui/form/Select/FormSelect';
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
          styles={{
            multiValue: (props, data) => {
              const colorScheme = (() => getTagColorScheme(data.data.label))();

              return {
                ...props,
                fontSize: 'xs',
                fontWeight: 'normal',
                color: colors[colorScheme][700],
                border: '1px solid',
                borderColor: colors[colorScheme][200],
                backgroundColor: colors[colorScheme][50],

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
      <span
        className='hover:border-gray-300 border-b border-transparent cursor-text text-gray-400 transition-colors duration-200 ease-in-out'
        onClick={(e) => {
          if (isCardOpen) {
            e.stopPropagation();
          }
          setIsFocused(true);
        }}
      >
        {placeholder}
      </span>
    );
  }

  return (
    <div
      className='flex gap-1 mt-2 pb-2 flex-wrap'
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
    </div>
  );
};

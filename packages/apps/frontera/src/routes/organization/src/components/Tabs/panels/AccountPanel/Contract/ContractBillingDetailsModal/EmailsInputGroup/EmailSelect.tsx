import { forwardRef } from 'react';
import { SelectInstance } from 'react-select';

import { SelectOption } from '@shared/types/SelectOptions.ts';

import { EmailMultiCreatableSelect } from './EmailMultiCreatableSelect';

interface EmailParticipantSelect {
  value: string[];
  isMulti: boolean;
  entryType: string;
  autofocus?: boolean;
  placeholder?: string;
  onChange: (value: SelectOption<string>[]) => void;
}

export const EmailSelect = forwardRef<SelectInstance, EmailParticipantSelect>(
  (
    { entryType, isMulti, placeholder = 'Enter email', value, onChange },
    ref,
  ) => {
    return (
      <div className='text-base group'>
        <label className='font-semibold text-sm'>{entryType}</label>
        <EmailMultiCreatableSelect
          ref={ref}
          value={value?.map((e) => ({ label: e, value: e }))}
          onChange={onChange}
          isMulti={isMulti}
          placeholder={placeholder}
          navigateAfterAddingToPeople={true}
          noOptionsMessage={() => null}
          // @ts-expect-error fix later
          getOptionLabel={(d) => {
            if (d?.__isNew__ || d.label === d.value) {
              return `${d.label}`;
            }

            return `${d.label} - ${d.value}`;
          }}
        />
      </div>
    );
  },
);

import { forwardRef } from 'react';
import { SelectInstance } from 'react-select';

import { SelectOption } from '@shared/types/SelectOptions.ts';

import { EmailMultiCreatableSelect } from './EmailMultiCreatableSelect';

interface EmailParticipantSelect {
  value: string[];
  entryType: string;
  autofocus?: boolean;
  placeholder?: string;
  onChange: (value: SelectOption<string>[]) => void;
}

export const EmailSelect = forwardRef<SelectInstance, EmailParticipantSelect>(
  ({ entryType, placeholder = 'Enter email', value, onChange }, ref) => {
    return (
      <div className='text-sm'>
        <label className='font-semibold text-sm'>{entryType}</label>
        <EmailMultiCreatableSelect
          ref={ref}
          value={value?.map((e) => ({ label: e, value: e }))}
          onChange={onChange}
          placeholder={placeholder}
          navigateAfterAddingToPeople={true}
          noOptionsMessage={() => null}
          // @ts-expect-error fix later
          getOptionLabel={(d) => {
            if (d?.__isNew__) {
              return `${d.label}`;
            }

            return `${d.label} - ${d.value}`;
          }}
        />
      </div>
    );
  },
);

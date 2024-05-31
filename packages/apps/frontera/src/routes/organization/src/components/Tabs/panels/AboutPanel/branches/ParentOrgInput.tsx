import React, { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Select } from '@ui/form/Select';
import { useStore } from '@shared/hooks/useStore';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';

interface ParentOrgInputProps {
  id: string;
  isReadOnly?: boolean;
  parentOrg: { label: string; value: string } | null;
}

export const ParentOrgInput = observer(
  ({ id, parentOrg, isReadOnly }: ParentOrgInputProps) => {
    const [searchTerm, setSearchTerm] = useState('');

    const store = useStore();
    const data = store.organizations.toComputedArray((arr) => {
      if (searchTerm) {
        arr = arr.filter((org) =>
          org.value.name.toLowerCase().includes(searchTerm.toLowerCase()),
        );
      }

      return arr;
    });

    const organization = store.organizations.value.get(id);

    const options = data
      .filter(
        (e) =>
          !e.value.subsidiaries?.length &&
          e.value.metadata?.id !== id &&
          e.value.name?.length > 0,
      )
      .map((org) => ({
        value: org.value.metadata?.id,
        label: org.value.name,
      }));

    return (
      <Select
        isClearable
        isReadOnly={isReadOnly}
        value={
          organization?.value?.parentCompanies?.length
            ? organization?.value?.parentCompanies.map((org) => ({
                value: org.organization?.metadata?.id,
                label: org.organization?.name,
              }))
            : ''
        }
        onChange={(e) => {
          organization?.update((org) => {
            if (org.parentCompanies.length === 0) {
              org.parentCompanies = [e.value];
            } else {
              org.parentCompanies = [];
            }

            return org;
          });
        }}
        onInputChange={(inputValue) => setSearchTerm(inputValue)}
        options={options || []}
        placeholder='Parent organization'
        leftElement={<ArrowCircleBrokenUpLeft className='text-gray-500 mr-3' />}
      />
    );
  },
);

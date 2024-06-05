import { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Select } from '@ui/form/Select';
import { useStore } from '@shared/hooks/useStore';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';

interface ParentOrgInputProps {
  id: string;
  isReadOnly?: boolean;
}

export const ParentOrgInput = observer(
  ({ id, isReadOnly }: ParentOrgInputProps) => {
    const store = useStore();
    const selectRef = useRef(null);
    const data = store.organizations?.toComputedArray((arr) => {
      return arr;
    });

    const organization = store.organizations.value.get(id);

    const options = data
      ?.filter((e) => e.value.metadata?.id !== id && e.value.name?.length > 0)
      .map((org) => ({
        value: org.value.metadata?.id,
        label: org.value.name,
      }));

    const selection = organization
      ? {
          value:
            organization?.value.parentCompanies[0]?.organization?.metadata?.id,
          label: organization?.value.parentCompanies[0]?.organization?.name,
        }
      : { value: '', label: '' };

    return (
      <Select
        isClearable
        ref={selectRef}
        isReadOnly={isReadOnly}
        value={selection}
        onChange={(e) => {
          const findOrg = store.organizations.value.get(e?.value);

          if (!e) {
            organization?.update((org) => {
              org.parentCompanies = [];

              return org;
            });

            return;
          }
          if (!findOrg) return;

          findOrg?.update((org) => {
            if (!organization) return org;
            org.subsidiaries = [{ organization: organization?.value }];

            return org;
          });
        }}
        options={options || []}
        placeholder='Parent organization'
        leftElement={<ArrowCircleBrokenUpLeft className='text-gray-500 mr-3' />}
      />
    );
  },
);

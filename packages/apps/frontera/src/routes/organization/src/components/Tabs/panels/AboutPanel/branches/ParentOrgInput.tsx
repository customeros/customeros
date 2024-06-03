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

    const selection =
      organization && organization?.value.parentCompanies.length === 1
        ? {
            value:
              organization?.value.parentCompanies[0].organization.metadata.id,
            label: organization?.value.parentCompanies[0].organization.name,
          }
        : null;

    return (
      <Select
        isClearable
        isReadOnly={isReadOnly}
        value={selection}
        onChange={(e) => {
          organization?.update((org) => {
            if (org.parentCompanies.length === 0) {
              org.parentCompanies.push({
                organization: {
                  metadata: { id: e?.value },
                  name: e?.label,
                  // eslint-disable-next-line @typescript-eslint/no-explicit-any
                } as any,
              });
            } else {
              org.parentCompanies = [];
            }

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

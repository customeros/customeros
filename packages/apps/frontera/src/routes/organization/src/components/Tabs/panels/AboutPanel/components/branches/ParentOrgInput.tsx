import { observer } from 'mobx-react-lite';

import { Select } from '@ui/form/Select';
import { Organization } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';

interface ParentOrgInputProps {
  id: string;
  isReadOnly?: boolean;
}

export const ParentOrgInput = observer(
  ({ id, isReadOnly }: ParentOrgInputProps) => {
    const store = useStore();
    const data = store.organizations?.toArray();

    const organization = store.organizations.value.get(id);

    const options = data
      ?.filter((e) => e.value.metadata?.id !== id && e.value.name?.length > 0)
      .map((org) => ({
        value: org.value.metadata?.id,
        label: org.value.name,
      }));

    const selection = organization
      ? {
          value: organization?.parentCompanies?.[0]?.metadata?.id,
          label: organization?.parentCompanies?.[0]?.name,
        }
      : { value: '', label: '' };

    return (
      <Select
        isClearable
        isReadOnly={isReadOnly}
        options={options || []}
        placeholder='Parent organization'
        value={selection.label ? selection : null}
        leftElement={<ArrowCircleBrokenUpLeft className='text-gray-500 mr-3' />}
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

            org.subsidiaries = [
              {
                organization: {
                  id: organization?.value?.metadata?.id,
                  name: organization?.value?.name,
                  metadata: { ...organization?.value?.metadata },
                } as Organization,
              },
            ];

            return org;
          });
        }}
      />
    );
  },
);

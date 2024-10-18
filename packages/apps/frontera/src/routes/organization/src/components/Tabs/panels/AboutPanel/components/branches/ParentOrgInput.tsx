import { runInAction } from 'mobx';
import { observer } from 'mobx-react-lite';

import { Combobox } from '@ui/form/Combobox';
import { Organization } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';

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

    const parentCompany = organization?.parentCompanies?.[0];

    const selection = parentCompany
      ? {
          value: parentCompany?.metadata?.id,
          label: parentCompany?.name,
        }
      : null;

    return (
      <Popover open={isReadOnly ? false : undefined}>
        <PopoverTrigger asChild className='cursor-pointer'>
          <div className='flex items-center min-h-10'>
            <ArrowCircleBrokenUpLeft className='text-gray-500 mr-3' />
            {parentCompany ? (
              <span>{selection?.label}</span>
            ) : (
              <span className='text-gray-400'>Parent organization</span>
            )}
          </div>
        </PopoverTrigger>
        <PopoverContent align='start' className='min-w-[264px] max-w-[320px]'>
          <Combobox
            isClearable
            value={selection}
            isReadOnly={isReadOnly}
            options={options || []}
            placeholder='Search...'
            onChange={(option) => {
              runInAction(() => {
                if (!organization) return;

                const newParent = store.organizations.value?.get(option?.value);

                if (!newParent) {
                  const parentId =
                    organization.value.parentCompanies?.[0]?.organization
                      ?.metadata?.id;

                  organization.value.parentCompanies = [];
                  organization.commit();

                  const parentCompany =
                    store.organizations.value?.get(parentId);

                  if (!parentCompany) return;

                  parentCompany.value.subsidiaries =
                    parentCompany.value.subsidiaries.filter(
                      (s) => s.organization.metadata.id !== organization.id,
                    );

                  parentCompany.commit();
                } else {
                  const currentParentId =
                    organization.value.parentCompanies?.[0]?.organization
                      ?.metadata?.id;

                  const currentParent =
                    store.organizations.value.get(currentParentId);

                  if (currentParent) {
                    const subsidiaryIndex =
                      currentParent.value.subsidiaries.findIndex(
                        (s) => s.organization.metadata.id === organization.id,
                      );

                    currentParent.value.subsidiaries.splice(subsidiaryIndex, 1);

                    currentParent.commit();

                    organization.value.parentCompanies = [];
                    organization.commit();
                  }

                  newParent.value?.subsidiaries?.push({
                    organization: {
                      id: organization?.value?.metadata?.id,
                      name: organization?.value?.name,
                      metadata: { ...organization?.value?.metadata },
                    } as Organization,
                  });
                  newParent.commit();

                  if (!Array.isArray(!organization.value.parentCompanies)) {
                    organization.value.parentCompanies = [];
                  }

                  organization.value.parentCompanies[0] = {
                    organization: newParent.value,
                  };

                  organization.commit();
                }
              });
            }}
          />
        </PopoverContent>
      </Popover>
    );
  },
);

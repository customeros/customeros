import { runInAction } from 'mobx';
import { observer } from 'mobx-react-lite';

import { Combobox } from '@ui/form/Combobox';
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
    const data = store.organizations?.findMany(
      (org) => org.id !== id && org.value.name.length > 0,
    );

    const organization = store.organizations.getById(id);

    const options = data.map((org) => ({
      value: org.value.metadata?.id,
      label: org.value.name,
    }));

    const parentCompany =
      organization?.value?.parentCompanies?.[0]?.organization;

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

                const newParent = store.organizations.getById(option?.value);

                if (!newParent) {
                  const parentId =
                    organization.value?.parentCompanies?.[0]?.organization
                      ?.metadata?.id;

                  organization.draft();
                  organization.clearParentCompanies();
                  organization.commit();

                  const parentCompany = store.organizations.getById(parentId!);

                  if (!parentCompany) return;

                  parentCompany.draft();
                  parentCompany.removeSubsidiary(organization.id);
                  parentCompany.commit();
                } else {
                  const currentParentId =
                    organization.value?.parentCompanies?.[0]?.organization
                      ?.metadata?.id;

                  const currentParent = store.organizations.getById(
                    currentParentId!,
                  );

                  if (currentParent) {
                    currentParent.draft();
                    currentParent.removeSubsidiary(organization.id);
                    currentParent.commit();

                    organization.draft();
                    organization.clearParentCompanies();
                    organization.commit();
                  }

                  newParent.draft();
                  newParent.addSubsidiary(organization.id);
                  newParent.commit();

                  organization.draft();

                  if (!Array.isArray(!organization.value?.parentCompanies)) {
                    organization.clearParentCompanies();
                  }

                  if (newParent?.value) {
                    organization.addParent(newParent.id);
                    organization.commit();
                  }
                }
              });
            }}
          />
        </PopoverContent>
      </Popover>
    );
  },
);

import React, { useState } from 'react';
import { useCreateOrganization } from '@spaces/hooks/useOrganization';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import {
  useAddOrganizationSubsidiary,
  useOrganizationSubsidiaries,
} from '@spaces/hooks/useOrganizationSubsidiaries';
import { DebouncedAutocomplete } from '@spaces/atoms/autocomplete';
import { useOrganizationSuggestionsList } from '@spaces/hooks/useOrganizations';
import { OrganizationSubsidiary } from '@spaces/organization/organization-details/subsidiaries/OrganizationSubsidiary';
import styles from './organization-subsidiaries.module.scss';
import PlusCircle from '@spaces/atoms/icons/PlusCircle';
import { OrganizationSubsidiariesSkeleton } from '@spaces/organization/organization-details/subsidiaries/skeletons';
export const OrganizationSubsidiaries = ({ id }: { id: string }) => {
  const { getOrganizationSuggestions } = useOrganizationSuggestionsList();
  const { data, loading, error } = useOrganizationSubsidiaries({ id });
  const { onCreateOrganization } = useCreateOrganization();
  const [organizationOptions, setOrganizationOptions] = useState<
    Array<{ value: string; label: string }>
  >([]);
  const [{ isEditMode }] = useRecoilState(organizationDetailsEdit);

  const { onAddOrganizationSubsidiary } = useAddOrganizationSubsidiary({ id });

  if (loading) {
    return <OrganizationSubsidiariesSkeleton />;
  }

  return (
    <article className={styles.subsidiary_section}>
      <h1 className={styles.subsidiary_header}>Branches</h1>

        <OrganizationSubsidiary
          subsidiaries={data?.subsidiaries || []}
          id={id}
        />

      {isEditMode && (
        <div className={styles.subsidiary_input}>
          <PlusCircle height={14} />
          <DebouncedAutocomplete
            key={`${data?.subsidiaries.length}-subsidiary-organization-id`}
            mode='fit-content'
            editable={true}
            value={''}
            suggestions={organizationOptions}
            onChange={(e) =>
              onAddOrganizationSubsidiary({
                organizationId: id,
                subOrganizationId: e.value,
              })
            }
            onAddNew={(e) => onCreateOrganization({ name: e.value })}
            onSearch={(filter: string) =>
              getOrganizationSuggestions(filter).then((options) =>
                setOrganizationOptions(options),
              )
            }
            newItemLabel='name'
            placeholder='Organization'
          />
        </div>
      )}
    </article>
  );
};

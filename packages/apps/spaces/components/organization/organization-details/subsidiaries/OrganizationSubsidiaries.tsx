import React from 'react';
import { useCreateOrganization } from '@spaces/hooks/useOrganization';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import { useAddOrganizationSubsidiary } from '@spaces/hooks/useOrganizationSubsidiaries';
import { useOrganizationSuggestionsList } from '@spaces/hooks/useOrganizations';
import { OrganizationSubsidiary } from '@spaces/organization/organization-details/subsidiaries/OrganizationSubsidiary';
import styles from './organization-subsidiaries.module.scss';
import PlusCircle from '@spaces/atoms/icons/PlusCircle';
import {
  Select,
  CreatableSelectMenu,
  SelectInput,
  SelectWrapper,
} from '@spaces/ui/form/select';

export const OrganizationSubsidiaries = ({
  id,
  subsidiaries,
}: {
  id: string;
  subsidiaries?: Array<any>;
}) => {
  const { getOrganizationSuggestions, organizationSuggestions } =
    useOrganizationSuggestionsList();
  const { onCreateOrganization } = useCreateOrganization();
  const [{ isEditMode }] = useRecoilState(organizationDetailsEdit);
  const { onAddOrganizationSubsidiary, saving } = useAddOrganizationSubsidiary({
    id,
  });

  return (
    <article className={styles.subsidiary_section}>
      <h1 className={styles.subsidiary_header}>Branches</h1>

      <OrganizationSubsidiary subsidiaries={subsidiaries || []} id={id} />

      {isEditMode && (
        <div
          className={styles.subsidiary_input}
          key={`data-select-subsidiary-${subsidiaries?.length || 0}`}
        >
          <PlusCircle height={14} />
          <Select<string>
            onSelect={(val) => {
              return onAddOrganizationSubsidiary({
                organizationId: id,
                subOrganizationId: val,
              });
            }}
            onCreateNewOption={(val) =>
              onCreateOrganization({ name: val }).then((d?: any) => {
                return onAddOrganizationSubsidiary({
                  organizationId: id,
                  subOrganizationId: d?.id,
                });
              })
            }
            onChange={(filter) => getOrganizationSuggestions(filter)}
            options={organizationSuggestions}
          >
            <SelectWrapper>
              <SelectInput
                placeholder='Organization'
                readOnly={!isEditMode}
                saving={saving}
              />
              <CreatableSelectMenu />
            </SelectWrapper>
          </Select>
        </div>
      )}
    </article>
  );
};

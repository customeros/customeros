import React from 'react';
import styles from './job-roles-input.module.scss';
import { AddIconButton } from '@spaces/atoms/icon-button/AddIconButton';
import { DeleteIconButton } from '@spaces/atoms/icon-button/DeleteIconButton';
import { Checkbox } from '@spaces/atoms/checkbox';
import { EditableContentInput } from '@spaces/atoms/input';
import {
  useCreateContactJobRole,
  useRemoveJobRoleFromContactJobRole,
  useUpdateContactJobRole,
} from '@spaces/hooks/useContactJobRole';
import { useOrganizationSuggestionsList } from '@spaces/hooks/useOrganizations/useOrganizationSuggestionsList';
import { useCreateOrganization } from '@spaces/hooks/useOrganization/useCreateOrganization';
import classNames from 'classnames';
import {
  Select,
  CreatableSelectMenu,
  SelectInput,
  SelectWrapper,
} from '@spaces/ui/form/select';

interface JobRoleInputProps {
  contactId: string;
  organization: {
    id: string;
    name: string;
  };
  jobRole: string;
  roleId: string;
  isEditMode?: boolean;
  showAddButton?: boolean;
  primary: boolean;
}

export const JobRoleInput: React.FC<JobRoleInputProps> = ({
  contactId,
  roleId,
  organization,
  jobRole,
  primary,
  isEditMode,
  showAddButton = false,
}) => {
  const { onCreateContactJobRole } = useCreateContactJobRole({ contactId });
  const { getOrganizationSuggestions, organizationSuggestions, loading } =
    useOrganizationSuggestionsList();
  const { onCreateOrganization, saving } = useCreateOrganization();
  const { onUpdateContactJobRole } = useUpdateContactJobRole({ contactId });
  const { onRemoveContactJobRole } = useRemoveJobRoleFromContactJobRole({
    contactId,
  });
  if (!isEditMode && !organization?.id) {
    return <div />;
  }

  return (
    <div>
      <div
        className={classNames(styles.jobAndOrganizationInputs, {
          [styles.primary]: primary && !isEditMode,
        })}
      >
        {isEditMode && (
          <DeleteIconButton
            style={{ position: 'absolute', left: -16, top: 6 }}
            onDelete={() => onRemoveContactJobRole(roleId)}
          />
        )}
        {(isEditMode || !!jobRole?.length) && (
          <EditableContentInput
            id={`conatct-personal-details-last-name-job-role-${contactId}=${roleId}`}
            label='Job title'
            isEditMode={isEditMode}
            value={jobRole || ''}
            placeholder='Job title'
            onChange={(value: string) => {
              roleId
                ? onUpdateContactJobRole({
                    id: roleId,
                    jobTitle: value,
                    organizationId: organization?.id,
                    primary,
                  })
                : onCreateContactJobRole({
                    jobTitle: value,
                  });
            }}
          />
        )}

        {(isEditMode || !!organization?.name?.length) && (
          <div style={{ marginRight: 8 }}>
            <Select<string>
              onSelect={(val) => {
                roleId
                  ? onUpdateContactJobRole({
                      id: roleId,
                      jobTitle: jobRole,
                      organizationId: val,
                      primary,
                    })
                  : onCreateContactJobRole({ organizationId: val });
              }}
              onCreateNewOption={(val) => {
                onCreateOrganization({ name: val }).then((d) => {
                  if (!d?.id) return;
                  return roleId
                    ? onUpdateContactJobRole({
                        id: roleId,
                        jobTitle: jobRole,
                        organizationId: d?.id,
                        primary,
                      })
                    : onCreateContactJobRole({ organizationId: d?.id });
                });
              }}
              onChange={(filter) => getOrganizationSuggestions(filter)}
              value={organization?.id}
              options={
                organization?.id
                  ? [
                      ...organizationSuggestions,
                      { value: organization?.id, label: organization?.name },
                    ]
                  : organizationSuggestions
              }
            >
              <SelectWrapper>
                <SelectInput
                  saving={saving}
                  placeholder='Organization'
                  readOnly={!isEditMode}
                />
                {isEditMode && <CreatableSelectMenu />}
              </SelectWrapper>
            </Select>
          </div>
        )}

        {isEditMode && (
          <Checkbox
            type='radio'
            checked={primary}
            label='Primary'
            // @ts-expect-error revisit
            onChange={(e) => {
              roleId
                ? onUpdateContactJobRole({
                    id: roleId,
                    jobTitle: jobRole,
                    organizationId: organization?.id,
                    primary: !primary,
                  })
                : onCreateContactJobRole({ primary: !primary });
            }}
          />
        )}

        {showAddButton && isEditMode && (
          <AddIconButton
            style={{
              width: '24px',
              height: '16px',
              position: 'relative',
            }}
            onAdd={() => {
              onCreateContactJobRole({
                jobTitle: '',
                primary: false,
                organizationId: '',
              });
            }}
          />
        )}
      </div>
    </div>
  );
};

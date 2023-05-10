import React, { useState } from 'react';
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
import { DebouncedAutocomplete } from '@spaces/atoms/autocomplete';

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
  const [organizationOptions, setOrganizationOptions] = useState<
    Array<{ value: string; label: string }>
  >([]);
  const { onCreateContactJobRole } = useCreateContactJobRole({ contactId });
  const { getOrganizationSuggestions } = useOrganizationSuggestionsList();
  const { onCreateOrganization } = useCreateOrganization();
  const { onUpdateContactJobRole } = useUpdateContactJobRole({ contactId });
  const { onRemoveContactJobRole } = useRemoveJobRoleFromContactJobRole({
    contactId,
  });

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
          <DebouncedAutocomplete
            mode='fit-content'
            editable={isEditMode}
            value={organization?.name || ''}
            suggestions={organizationOptions}
            onChange={(e) =>
              roleId
                ? onUpdateContactJobRole({
                    id: roleId,
                    jobTitle: jobRole,
                    organizationId: e.value,
                    primary,
                  })
                : onCreateContactJobRole({ organizationId: e.value })
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

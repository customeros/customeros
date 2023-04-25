import React, { useEffect, useState } from 'react';
import styles from './job-roles-input.module.scss';
import {
  AddIconButton,
  Autocomplete,
  Checkbox,
  DeleteIconButton,
  EditableContentInput,
} from '../../../ui-kit/atoms';
import {
  useCreateContactJobRole,
  useRemoveJobRoleFromContactJobRole,
  useUpdateContactJobRole,
} from '../../../../hooks/useContactJobRole';
import { capitalizeFirstLetter } from '../../../../utils';
import { useOrganizationsOptions } from '../../../../hooks/useOrganizations';
import { useCreateOrganization } from '../../../../hooks/useOrganization';
import classNames from 'classnames';
import { useAutoAnimate } from '@formkit/auto-animate/react';

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
  const { data, loading, error } = useOrganizationsOptions();
  const { onCreateOrganization } = useCreateOrganization();
  const { onUpdateContactJobRole } = useUpdateContactJobRole({ contactId });
  const { onRemoveContactJobRole } = useRemoveJobRoleFromContactJobRole({
    contactId,
  });

  useEffect(() => {
    if (data) {
      const options = data.map(({ id, name }) => ({
        value: id,
        label: capitalizeFirstLetter(name),
      }));

      setOrganizationOptions(options);
    }
  }, [data]);
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
          <Autocomplete
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

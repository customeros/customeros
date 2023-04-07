import React, { useEffect, useState } from 'react';
import styles from './job-roles-input.module.scss';
import {
  Autocomplete,
  EditableContentInput,
  IconButton,
  Plus,
  Trash,
} from '../../../ui-kit/atoms';
import {
  useCreateContactJobRole,
  useRemoveJobRoleFromContactJobRole,
  useUpdateContactJobRole,
} from '../../../../hooks/useContactJobRole';
import { capitalizeFirstLetter } from '../../../../utils';
import { useOrganizationsOptions } from '../../../../hooks/useOrganizations';
import { useCreateOrganization } from '../../../../hooks/useOrganization';
import { useRouter } from 'next/router';

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
}

export const JobRoleInput: React.FC<JobRoleInputProps> = ({
  contactId,
  roleId,
  organization,
  jobRole,
  isEditMode,
  showAddButton = false,
}) => {
  const [organizationOptions, setOrganizationOptions] = useState<
    Array<{ value: string; label: string }>
  >([]);
  const { onCreateContactJobRole } = useCreateContactJobRole({ contactId });
  const { data, loading, error } = useOrganizationsOptions();
  const { onCreateOrganization } = useCreateOrganization();
  const {
    query: { id },
  } = useRouter();
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
      <div className={styles.jobAndOrganizationInputs}>
        {isEditMode && (
          <IconButton
            onClick={(e) => {
              e.stopPropagation();
              e.preventDefault();
              onRemoveContactJobRole(roleId);
            }}
            icon={<Trash style={{ transform: 'scale(0.7)' }} />}
            size='xxxxs'
            role='button'
            mode='text'
          />
        )}
        {(isEditMode || !!jobRole?.length) && (
          <EditableContentInput
            isEditMode={isEditMode}
            value={jobRole || ''}
            placeholder='Job title'
            onChange={(value: string) => {
              roleId
                ? onUpdateContactJobRole({
                    id: roleId,
                    jobTitle: value,
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
                    organizationId: e.value,
                  })
                : onCreateContactJobRole({ organizationId: e.value })
            }
            onAddNew={(e) => onCreateOrganization({ name: e.value })}
            newItemLabel='name'
            placeholder='Organization'
          />
        )}

        {showAddButton && isEditMode && (
          <IconButton
            onClick={(e) => {
              e.stopPropagation();
              e.preventDefault();
              onCreateContactJobRole({
                jobTitle: '',
                primary: false,
                organizationId: organization?.id,
              });
            }}
            icon={<Plus style={{ transform: 'scale(0.8)' }} />}
            size='xxxxs'
            role='button'
            mode='text'
          />
        )}
      </div>
    </div>
  );
};

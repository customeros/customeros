import React, { useEffect, useState } from 'react';
import styles from './job-roles-input.module.scss';
import {
  Autocomplete,
  DeleteIconButton,
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
import { useAttachOrganizationToContact } from '../../../../hooks/useContact';

interface JobRoleInputProps {
  contactId: string;
  organization?: {
    id: string;
    name: string;
  };
  isEditMode?: boolean;
  showAddNew?: boolean;
  index?: number;
}

export const AttachOrganizationInput: React.FC<JobRoleInputProps> = ({
  contactId,
  organization,
  isEditMode,
  showAddNew = false,
  index,
}) => {
  const [organizationOptions, setOrganizationOptions] = useState<
    Array<{ value: string; label: string }>
  >([]);
  const { onAttachOrganizationToContact } = useAttachOrganizationToContact({
    contactId,
  });
  const { data } = useOrganizationsOptions();
  const { onCreateOrganization } = useCreateOrganization();
  useEffect(() => {
    if (data) {
      const options = data.map(({ id, name }) => ({
        value: id,
        label: capitalizeFirstLetter(name),
      }));

      setOrganizationOptions(options);
    }
  }, [data]);

  if (!organization) {
    return (
      <div className={styles.attachOrgInput}>
        {isEditMode && (
          <div className={styles.addNewSection}>
            <Autocomplete
              mode='fit-content'
              editable={isEditMode}
              value={''}
              suggestions={organizationOptions}
              onChange={(e) =>
                onAttachOrganizationToContact({
                  contactId,
                  organizationId: e.value,
                })
              }
              onAddNew={(e) => onCreateOrganization({ name: e.value })}
              newItemLabel='name'
              placeholder='Organization'
            />
          </div>
        )}
      </div>
    );
  }

  return (
    <div className={styles.attachOrgInput}>
      {isEditMode && (
        <DeleteIconButton
          style={{ position: 'absolute', left: -16, top: 6 }}
          onDelete={() => null}
        />
      )}

      {(isEditMode || !!organization?.name?.length) && (
        <Autocomplete
          mode='fit-content'
          editable={isEditMode}
          value={organization?.name || ''}
          suggestions={organizationOptions}
          onChange={(e) =>
            onAttachOrganizationToContact({
              contactId,
              organizationId: e.value,
            })
          }
          onAddNew={(e) => onCreateOrganization({ name: e.value })}
          newItemLabel='name'
          placeholder='Organization'
        />
      )}
      {isEditMode && showAddNew && (
        <div className={styles.addNewSection}>
          <Autocomplete
            mode='fit-content'
            editable={isEditMode}
            value={''}
            suggestions={organizationOptions}
            onChange={(e) =>
              onAttachOrganizationToContact({
                contactId,
                organizationId: e.value,
              })
            }
            onAddNew={(e) => onCreateOrganization({ name: e.value })}
            newItemLabel='name'
            placeholder='Organization'
          />
        </div>
      )}
    </div>
  );
};

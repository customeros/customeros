import React, { useEffect, useState } from 'react';
import styles from './job-roles-input.module.scss';
import { Autocomplete, DeleteIconButton } from '../../../ui-kit';
import { capitalizeFirstLetter } from '../../../../utils';
import { useOrganizationsOptions } from '../../../../hooks/useOrganizations';
import { useCreateOrganization } from '../../../../hooks/useOrganization';
import { useAttachOrganizationToContact } from '../../../../hooks/useContact';
import { useRemoveOrganizationFromContact } from '../../../../hooks/useContact/useRemoveOrganizationFromContact';

interface JobRoleInputProps {
  contactId: string;
  organization?: {
    id: string;
    name: string;
  };
  isEditMode?: boolean;
  showAddNew?: boolean;
}

export const AttachOrganizationInput: React.FC<JobRoleInputProps> = ({
  contactId,
  organization,
  isEditMode,
  showAddNew = false,
}) => {
  const [organizationOptions, setOrganizationOptions] = useState<
    Array<{ value: string; label: string }>
  >([]);
  const { onAttachOrganizationToContact } = useAttachOrganizationToContact({
    contactId,
  });
  const { onRemoveOrganizationFromContact } = useRemoveOrganizationFromContact({
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
          onDelete={() =>
            onRemoveOrganizationFromContact({
              contactId,
              organizationId: organization?.id,
            })
          }
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

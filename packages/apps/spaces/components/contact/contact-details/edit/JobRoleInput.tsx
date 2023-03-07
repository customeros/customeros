import React, { useEffect, useRef, useState } from 'react';

import styles from './contact-details-edit.module.scss';
import {
  DebouncedInput,
  IconButton,
  Input,
  Trash,
} from '../../../ui-kit/atoms';
import {
  useCreateContactJobRole,
  useRemoveJobRoleFromContactJobRole,
  useUpdateContactJobRole,
} from '../../../../hooks/useContactJobRole';
import { Controller } from 'react-hook-form';
import { Dropdown } from 'primereact/dropdown';
import { OverlayPanel } from '../../../ui-kit/atoms/overlay-panel';
import { capitalizeFirstLetter } from '../../../../utils';
import { useOrganizationsOptions } from '../../../../hooks/useOrganizations';

interface JobRoleInputProps {
  contactId: string;
  roleId?: string;
  field: any;
  fields: any;
  index: number;
  register: any;
  control: any;
  append: any;
  remove: any;
}

export const JobRoleInput: React.FC<JobRoleInputProps> = ({
  contactId,
  roleId,
  index,
  register,
  fields,
  control,
  append,
  remove,
}) => {
  const organizationSelectorRef = useRef(null);

  const [organizationOptions, setOrganizationOptions] = useState<
    Array<{ value: string; label: string }>
  >([]);
  const { onCreateContactJobRole } = useCreateContactJobRole({ contactId });
  const { data, loading, error } = useOrganizationsOptions();
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
        <Controller
          control={control}
          render={({ field }) => (
            <DebouncedInput
              id={`contact-${contactId}-job-title`}
              // hideLabel={true}
              // label={'role'}
              className={styles.jobRoleInput}
              placeholder='Job title'
              inputSize='xxxs'
              onChange={(e) => {
                console.log('🏷️ ----- e: ', e, roleId);

                roleId
                  ? onUpdateContactJobRole({
                      id: roleId,
                      jobTitle: e.target.value,
                    })
                  : field.onChange(e.target.value);
              }}
            />
          )}
          name={`jobRoles.${index}.jobTitle`}
        />

        <span className={styles.copy}>at</span>

        <Controller
          control={control}
          name={`jobRoles.${index}.organizationId`}
          render={({ field }) => (
            <Dropdown
              id={field.name}
              value={field.value}
              onChange={(e) =>
                roleId
                  ? onUpdateContactJobRole({
                      id: roleId,
                      organizationId: e.value,
                    })
                  : field.onChange(e.value)
              }
              options={organizationOptions}
              optionValue='value'
              optionLabel='label'
              placeholder='Organization'
              className={styles.titleSelector}
            />
          )}
        />
        <OverlayPanel
          ref={organizationSelectorRef}
          model={organizationOptions}
        />
        <IconButton
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
            roleId ? onRemoveContactJobRole(roleId) : remove(index);
          }}
          icon={<Trash style={{ transform: 'scale(0.8)' }} />}
          size='xxxs'
          role='button'
          mode='text'
        />
      </div>
    </div>
  );
};

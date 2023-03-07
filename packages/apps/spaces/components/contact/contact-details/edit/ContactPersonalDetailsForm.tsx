import React, { FormEvent, useEffect, useMemo, useRef } from 'react';
import Image from 'next/image';
import { Button, User, UserPlus } from '../../../ui-kit/atoms';
import styles from '../contact-details.module.scss';
import { ContactTags } from '../../contact-tags';
import { ControlledInput } from '../../../ui-kit/atoms/input/ControlledInput';
import { Controller, useFieldArray, useForm } from 'react-hook-form';
import {
  JobRole,
  PersonTitle,
} from '../../../../graphQL/__generated__/generated';
import { OverlayPanel } from '../../../ui-kit/atoms/overlay-panel';
import { Dropdown } from 'primereact/dropdown';
import { capitalizeFirstLetter } from '../../../../utils';
import { JobRoleInput } from './JobRoleInput';
import { useCreateContactJobRole } from '../../../../hooks/useContactJobRole';

export const ContactPersonalDetailsEditForm = ({
  data,
  onSubmit,
  mode = 'EDIT',
}: {
  data: any;
  onSubmit: any;
  mode?: 'CREATE' | 'EDIT';
}) => {
  const titleSelectorRef = useRef(null);
  const { onCreateContactJobRole } = useCreateContactJobRole({
    contactId: data?.id,
  });

  const { control, getValues, register, reset } = useForm({
    defaultValues: {
      firstName: data?.firstName || '',
      id: data.id,
      label: data?.label || '',
      lastName: data?.lastName || '',
      title: data?.title || PersonTitle.Mr,
      jobRoles: data?.jobRoles || [],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'jobRoles',
  });

  useEffect(() => {
    if (data?.id) {
      reset({
        firstName: data?.firstName || '',
        id: data.id,
        label: data?.label || '',
        lastName: data?.lastName || '',
        title: data?.title || PersonTitle.Mr,
        jobRoles:
          data?.jobRoles.map((role: JobRole) => ({
            roleId: role.id,
            organizationId: role?.organization?.id,
            jobTitle: role.jobTitle || '',
          })) || [],
      });
    }
  }, []);

  const handleJobRolesUpdate = async (roles: Array<JobRole>) => {
    // now only 1 title works
    await onCreateContactJobRole(roles[0]);
  };

  const handleSubmit = (
    e: FormEvent<HTMLButtonElement>,
    { jobRoles, ...values }: any,
  ) => {
    e.stopPropagation();
    e.preventDefault();

    if (data.id) {
      // only in edit mode
      handleJobRolesUpdate(jobRoles);
    }

    onSubmit(values);
  };

  const titleOptions = useMemo(
    () =>
      Object.values(PersonTitle).map((labelOption) => ({
        value: labelOption,
        label: capitalizeFirstLetter(labelOption),
      })),
    [],
  );

  return (
    <form className={styles.header}>
      {mode === 'EDIT' && (
        <div className={styles.photo}>
          {data?.photo ? (
            <Image src={data?.photo} alt={''} height={40} width={40} />
          ) : (
            <User />
          )}
        </div>
      )}

      <div className={styles.name}>
        {mode === 'CREATE' && (
          <div style={{ display: 'flex', justifyContent: 'center' }}>
            <div className={styles.photo}>
              <UserPlus />
            </div>
          </div>
        )}

        <div
          className={styles.formFields}
          style={{ display: 'flex', flexDirection: 'column', width: '100%' }}
        >
          <label className={styles.titleSelectorLabel} htmlFor='title'>
            Title
          </label>
          <Controller
            name='title'
            control={control}
            render={({ field }) => (
              <Dropdown
                id={field.name}
                value={field.value}
                onChange={(e) => field.onChange(e.value)}
                options={titleOptions}
                optionValue='value'
                optionLabel='label'
                className={styles.titleSelector}
              />
            )}
          />
          <OverlayPanel ref={titleSelectorRef} model={titleOptions} />

          <ControlledInput
            id='first-name-contact-input'
            required={true}
            inputSize='xxxs'
            control={control}
            name='firstName'
            placeholder='First name'
            label='First name'
            hideLabel
          />
          <ControlledInput
            id='first-name-contact-input'
            required={true}
            inputSize='xxxs'
            control={control}
            name='lastName'
            placeholder='Last name'
            label='Last name'
            hideLabel
          />
        </div>

        {mode === 'EDIT' &&
          fields.map((field, index) => {
            return (
              <JobRoleInput
                append={append}
                remove={remove}
                contactId={data.id}
                //@ts-expect-error fixme
                roleId={field?.roleId}
                key={field.id}
                field={field}
                fields={fields}
                control={control}
                index={index}
                register={register}
              />
            );
          })}

        {mode === 'EDIT' && fields.length < 1 && (
          <div className={styles.buttonWrapper}>
            <Button
              className={styles.addNewJobRoleButton}
              onClick={() => append({ jobTitle: '', organizationId: '' })}
              mode='link'
              type='button'
              icon={<span>+</span>}
            >
              Add job role
            </Button>
          </div>
        )}

        {mode === 'EDIT' && <ContactTags id={data.id} mode={'EDIT'} />}

        <div
          style={{
            display: 'flex',
            justifyContent: 'flex-end',
            marginTop: '8px',
          }}
        >
          <Button
            mode='primary'
            type='button'
            onClick={(e) => handleSubmit(e, getValues())}
          >
            {mode === 'CREATE' ? 'Create contact' : 'Save'}
          </Button>
        </div>
      </div>
    </form>
  );
};

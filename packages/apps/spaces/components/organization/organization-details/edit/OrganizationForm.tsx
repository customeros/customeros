import React, { FormEvent } from 'react';

import styles from '../organization-details.module.scss';
import { useForm } from 'react-hook-form';
import { Organization } from '../../../../graphQL/__generated__/generated';
import { ControlledInput } from '../../../ui-kit/atoms/input/ControlledInput';
import { Button } from '../../../ui-kit';
import { ControlledTextArea } from '../../../ui-kit/atoms/input/ControlledTextarea';

interface OrganizationFormProps {
  data?: Organization;
  onSubmit: any;
  mode?: 'EDIT' | 'CREATE';
}

export const OrganizationForm: React.FC<OrganizationFormProps> = ({
  data,
  onSubmit,
  mode = 'EDIT',
}) => {
  const { control, reset, setValue, getValues, register } = useForm({
    defaultValues: {
      name: data?.name || '',
      description: data?.description || '',
      website: data?.website || '',
      industry: data?.industry || '',
    },
  });

  const handleSubmit = (e: FormEvent<HTMLFormElement>, values: any) => {
    e.stopPropagation();
    e.preventDefault();
    onSubmit(values);
  };

  return (
    <form
      className={styles.organizationDetails}
      onSubmit={(e) => handleSubmit(e, getValues())}
    >
      <div className={styles.bg}>
        <ControlledInput
          required={true}
          inputSize='xxxs'
          control={control}
          name='name'
          placeholder=''
          label='Name'
        />
        <ControlledInput
          required={true}
          inputSize='xxxs'
          control={control}
          name='industry'
          placeholder=''
          label='Industry'
        />
        <ControlledTextArea
          required={true}
          inputSize='xxxs'
          control={control}
          name='description'
          placeholder=''
          label='Description'
        />

        <ControlledInput
          required={true}
          inputSize='xxxs'
          control={control}
          name='website'
          placeholder=''
          label='Website'
        />

        <div
          style={{
            display: 'flex',
            justifyContent: 'flex-end',
            marginTop: '8px',
          }}
        >
          <Button mode='primary' type='submit'>
            {mode === 'CREATE' ? 'Create organization' : 'Save'}
          </Button>
        </div>
      </div>
    </form>
  );
};

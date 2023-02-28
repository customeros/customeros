import React, { FormEvent } from 'react';

import styles from '../organization-details.module.scss';
import { useForm } from 'react-hook-form';
import { Organization } from '../../../../graphQL/__generated__/generated';

interface OrganizationFormProps {
  data: Organization;
  onSubmit: any;
}

export const OrganizationForm: React.FC<OrganizationFormProps> = ({
  data,
  onSubmit,
}) => {
  const { control, reset, setValue, getValues, register } = useForm({
    defaultValues: {
      name: '',
      description: '',
      domain: '',
      domains: [],
      industry: '',
      isPublic: '',
      organizationTypeId: '',
      appSource: '',
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
      TODO
    </form>
  );
};

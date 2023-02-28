import React, { FormEvent, useRef } from 'react';
import Image from 'next/image';
import { Button } from '../../ui-kit/atoms';
import styles from './contact-details.module.scss';
import { ContactTags } from '../contact-tags';
import { ControlledInput } from '../../ui-kit/atoms/input/ControlledInput';
import { Controller, useForm } from 'react-hook-form';
import { ContactOwner } from './ContactOwner';
import { useUpdateContactPersonalDetails } from '../../../hooks/useContact';
import { PersonTitle } from '../../../graphQL/__generated__/generated';
import { OverlayPanel } from '../../ui-kit/atoms/overlay-panel';
import { Dropdown } from 'primereact/dropdown';
import { capitalizeFirstLetter } from '../../../utils';

export const ContactPersonalDetailsEdit = ({
  data,
  onSetMode,
}: {
  data: any;
  onSetMode: any;
}) => {
  const titleSelectorRef = useRef(null);

  const { control, reset, setValue, getValues, register } = useForm({
    defaultValues: {
      firstName: data?.firstName || '',
      id: data.id,
      label: data?.label || '',
      lastName: data?.lastName || '',
      ownerId: data?.ownerId || '',
      title: data?.title || PersonTitle.Mr,
    },
  });
  const { onUpdateContactPersonalDetails } = useUpdateContactPersonalDetails({
    contactId: data.id,
  });

  const handleUpdateDetails = (e: FormEvent<HTMLFormElement>, values: any) => {
    e.stopPropagation();
    e.preventDefault();
    onUpdateContactPersonalDetails(values).then(() => {
      onSetMode('PREVIEW');
    });
  };

  const titleOptions = Object.values(PersonTitle).map((labelOption) => ({
    value: labelOption,
    label: capitalizeFirstLetter(labelOption),
  }));
  return (
    <form
      className={styles.header}
      onSubmit={(e) => handleUpdateDetails(e, getValues())}
    >
      <div className={styles.photo}>
        {data?.photo ? (
          <Image src={data?.photo} alt={''} height={40} width={40} />
        ) : (
          <div>{data?.firstName?.[0]}</div>
        )}
      </div>
      <div className={styles.name}>
        <div
          className={styles.formFields}
          style={{ display: 'flex', flexDirection: 'column', width: '100%' }}
        >
          <label className={styles.xyz} htmlFor='title'>
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
            required={true}
            inputSize='xxxs'
            control={control}
            name='firstName'
            placeholder='Type in first name'
            label='First name'
          />
          <ControlledInput
            required={true}
            inputSize='xxxs'
            control={control}
            name='lastName'
            placeholder='Type in last name'
            label='Last name'
          />
        </div>
        <ContactOwner control={control} setValue={setValue} />
        <ContactTags id={data.id} mode={'EDIT'} />
        <div
          style={{
            display: 'flex',
            justifyContent: 'flex-end',
            marginTop: '8px',
          }}
        >
          <Button mode='primary' type='submit'>
            Save
          </Button>
        </div>
      </div>
    </form>
  );
};

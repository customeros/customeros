import React, { FC } from 'react';
import { NoteEditor } from '../../ui-kit/molecules';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { useCreateOrganizationNote } from '../../../hooks/useOrganizationNote';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}
interface Props {
  mode: NoteEditorModes;
  organizationId: string;
}

const DEFAULT_VALUES = {
  html: '',
  htmlEnhanced: '',
};
export const OrganizationNoteEditor: FC<Props> = ({
  mode = NoteEditorModes.ADD,
  organizationId,
}) => {
  const { handleSubmit, setValue, getValues, control, reset } = useForm({
    defaultValues: DEFAULT_VALUES,
  });

  const { onCreateOrganizationNote } = useCreateOrganizationNote({
    organizationId,
  });
  const isEditMode = mode === NoteEditorModes.EDIT;

  const onSubmit = handleSubmit(async (d) => {
    //remove src attribute to not send the file bytes in here

    const dataToSubmit = {
      appSource: 'Openline',
      html: d?.htmlEnhanced?.replaceAll(/.src(\S*)/g, ''),
    };

    onCreateOrganizationNote(dataToSubmit).then(() => reset(DEFAULT_VALUES));

    toast.success('Note added successfully!');
  });

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: isEditMode ? '-17px -24px' : 0,
      }}
    >
      <Controller
        name='htmlEnhanced'
        control={control}
        render={({ field }) => (
          <NoteEditor
            mode={NoteEditorModes.ADD}
            onGetFieldValue={getValues}
            value={field.value}
            onSave={onSubmit}
            onTextChange={(e) => setValue('htmlEnhanced', e.htmlValue)}
          />
        )}
      />
    </div>
  );
};

import React, { ButtonHTMLAttributes, ChangeEvent, FC } from 'react';
import { Control, Controller, FieldValues, useForm } from 'react-hook-form';
import { NoteEditor } from '../../ui-kit/molecules/note-editor';
import { Button } from '../../ui-kit';
import { useCreateContact } from '../../../hooks/useContact';
import { useCreateContactNote } from '../../../hooks/useContactNote';
import { toast } from 'react-toastify';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}
interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  fieldName: string;
  fieldController: Control<{ id: any; html: any; htmlEnhanced: any }, any>;
  onGetFieldValue: any;
  mode: NoteEditorModes;
  onTextChange: (e: any) => void;
  onSave: () => void;
  onCancel?: () => void;
}

const DEFAULT_VALUES = {
  html: '',
  htmlEnhanced: '',
};
export const ContactNoteEditor: FC<any> = ({
  mode,
  contactId,
  onTextChange,
  onSave,
  onCancel = () => null,
}) => {
  const { register, handleSubmit, setValue, getValues, control, reset } =
    useForm({
      defaultValues: DEFAULT_VALUES,
    });

  const { onCreateContactNote } = useCreateContactNote({ contactId });
  const isEditMode = mode === NoteEditorModes.EDIT;

  const onSubmit = handleSubmit(async (d) => {
    //remove src attribute to not send the file bytes in here

    const dataToSubmit = {
      appSource: 'Openline',
      html: d?.htmlEnhanced?.replaceAll(/.src(\S*)/g, ''),
    };

    onCreateContactNote(dataToSubmit).then(() => reset(DEFAULT_VALUES));

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
            value={field.value}
            onSave={onSubmit}
            onTextChange={(e) => setValue('htmlEnhanced', e.htmlValue)}
          />
        )}
      />
    </div>
  );
};

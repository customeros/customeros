import React from 'react';
import { Controller, useForm } from 'react-hook-form';
import { useUpdateContactNote } from '../../../hooks/useContactNote';
import { NoteEditorModes } from './ContactEditor';
import { Editor } from '../../ui-kit/molecules';

interface Props {
  contactId: string;
  note: any;
  isEdit: boolean;
  onCancel: () => void;
  onSuccess: (data: any) => void;
}

function ContactNoteModalTemplate(props: Props) {
  const { onUpdateContactNote } = useUpdateContactNote();
  const { handleSubmit, setValue, getValues, control, reset } = useForm({
    defaultValues: {
      id: props.note?.id || '',
      html: props.note?.html || '',
      htmlEnhanced: props.note.htmlEnhanced || '',
    },
  });

  const onSubmit = handleSubmit(({ htmlEnhanced, ...data }) => {
    const dataToSubmit = {
      ...data,
      html: htmlEnhanced?.replaceAll(/.src(\S*)/g, ''), //remove src attribute to not send the file bytes in here
    };
    onUpdateContactNote(dataToSubmit).then(() => {
      props.onSuccess(dataToSubmit);
      reset({
        id: '',
        html: '',
        htmlEnhanced: '',
      });
    });
  });

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: props.isEdit ? '-17px -24px' : 0,
      }}
    >
      <Controller
        name='htmlEnhanced'
        control={control}
        render={({ field }) => (
          <Editor
            onCancel={() => props.onCancel()}
            mode={NoteEditorModes.EDIT}
            onGetFieldValue={getValues}
            value={field.value}
            onSave={onSubmit}
            label='Save'
            onTextChange={(e) => setValue('htmlEnhanced', e.htmlValue)}
          />
        )}
      />
    </div>
  );
}

export default ContactNoteModalTemplate;

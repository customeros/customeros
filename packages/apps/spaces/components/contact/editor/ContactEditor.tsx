import React, { FC } from 'react';
import { Editor } from '../../ui-kit/molecules';
import { Controller, useForm } from 'react-hook-form';
import { useCreateContactNote } from '../../../hooks/useNote';
import { useRecoilState, useRecoilValue } from 'recoil';
import { editorEmail, editorMode, EditorMode } from '../../../state';
import { EmailFields } from './email-fields';
import classNames from 'classnames';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}

interface Props {
  mode: NoteEditorModes;
  contactId: string;
}

const DEFAULT_VALUES = {
  html: '',
  htmlEnhanced: '',
};
export const ContactEditor: FC<Props> = ({
  mode = NoteEditorModes.ADD,
  contactId,
}) => {
  const { handleSubmit, setValue, getValues, control, reset } = useForm({
    defaultValues: DEFAULT_VALUES,
  });
  const [editorModeState, setMode] = useRecoilState(editorMode);
  const {
    handleSubmit: handleSendEmail,
    to,
    respondTo,
  } = useRecoilValue(editorEmail);
  const { onCreateContactNote, saving } = useCreateContactNote({ contactId });

  const onSubmit = handleSubmit(async (d) => {
    //remove src attribute to not send the file bytes in here
    // identiti header - from session uuid
    const dataToSubmit = {
      appSource: 'Openline',
      html: d?.htmlEnhanced?.replaceAll(/.src(\S*)/g, '') || '',
    };

    editorModeState.mode === EditorMode.Email && handleSendEmail
      ? handleSendEmail(
          dataToSubmit.html.replace(/(<([^>]+)>)/gi, ''),
          () => reset(DEFAULT_VALUES),
          to,
          respondTo,
        )
      : onCreateContactNote(dataToSubmit).then(() => reset(DEFAULT_VALUES));
  });

  const handleCancel = () => {
    setMode({ mode: EditorMode.Note, submitButtonLabel: 'Log into timeline' });
    reset(DEFAULT_VALUES);
  };

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: 0,
      }}
      key={editorModeState.mode}
    >
      <Controller
        name='htmlEnhanced'
        control={control}
        render={({ field }) => (
          <div
            className={classNames({
              'openline-editor-email':
                editorModeState.mode === EditorMode.Email,
            })}
          >
            {editorModeState.mode === EditorMode.Email && <EmailFields />}
            <Editor
              mode={NoteEditorModes.ADD}
              onGetFieldValue={getValues}
              value={field.value}
              saving={saving}
              onSave={onSubmit}
              onCancel={
                editorModeState.mode === EditorMode.Email
                  ? handleCancel
                  : undefined
              }
              label={editorModeState.submitButtonLabel}
              onHtmlChanged={(html: string) => setValue('htmlEnhanced', html)}
            />
          </div>
        )}
      />
    </div>
  );
};

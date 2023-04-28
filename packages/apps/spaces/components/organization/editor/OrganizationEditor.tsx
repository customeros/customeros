import React, { FC, useCallback, useRef } from 'react';
import classNames from 'classnames';
import { Controller, useForm } from 'react-hook-form';
import { useCreateOrganizationNote } from '../../../hooks/useNote';
import { editorEmail, editorMode, EditorMode, userData } from '../../../state';
import { EmailFields } from '../../contact/editor/email-fields';
import { useRecoilState, useRecoilValue } from 'recoil';
import {
  extraAttributes,
  SocialEditor,
} from '../../ui-kit/molecules/editor/SocialEditor';
import {
  AnnotationExtension,
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  FontSizeExtension,
  HistoryExtension,
  ImageExtension,
  ItalicExtension,
  LinkExtension,
  MentionAtomExtension,
  OrderedListExtension,
  StrikeExtension,
  TextColorExtension,
  UnderlineExtension,
  wysiwygPreset,
} from 'remirror/extensions';
import { toast } from 'react-toastify';
import { prosemirrorNodeToHtml } from 'remirror';
import { useRemirror } from '@remirror/react';
import { TableExtension } from '@remirror/extension-react-tables';
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
export const OrganizationEditor: FC<Props> = ({
  mode = NoteEditorModes.ADD,
  organizationId,
}) => {
  const [editorModeState, setMode] = useRecoilState(editorMode);

  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),

    ...wysiwygPreset(),
    new BoldExtension(),
    new ItalicExtension(),
    new BlockquoteExtension(),
    new ImageExtension({
      // uploadHandler: (e) => console.log('upload handler', e),
    }),
    new LinkExtension({ autoLink: true }),
    new TextColorExtension(),
    new UnderlineExtension(),
    new FontSizeExtension(),
    new HistoryExtension(),
    new AnnotationExtension(),
    new BulletListExtension(),
    new OrderedListExtension(),
    new StrikeExtension(),
  ];
  const extensions = useCallback(
    () => [...remirrorExtentions],
    [editorModeState.mode],
  );

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',
  });
  const { identity: loggedInUserEmail } = useRecoilValue(userData);
  const {
    handleSubmit: handleSendEmail,
    to,
    respondTo,
  } = useRecoilValue(editorEmail);

  const editorRef = useRef<any | null>(null);

  const { onCreateOrganizationNote, saving } = useCreateOrganizationNote({
    organizationId,
  });
  const isEditMode = mode === NoteEditorModes.EDIT;
  const handleResetEditor = (res: any) => {
    if (!res || !res?.id) return;
    const context = getContext();
    if (context) {
      context.commands.resetContent();
    }
  };
  const submitButtonOptions = [
    {
      label: 'Log as Note',
      command: () => {
        const data = prosemirrorNodeToHtml(state.doc);
        const dataToSubmit = {
          appSource: 'Openline',
          html: data?.replaceAll(/.src(\S*)/g, '') || '',
        };
        return onCreateOrganizationNote(dataToSubmit).then((res) =>
          handleResetEditor(res),
        );
      },
    },
  ];
  const submitEmailButtonOptions = [
    {
      label: 'Send Email',
      command: () => {
        const data = prosemirrorNodeToHtml(state.doc);
        if (!handleSendEmail) {
          toast.error('Client error occurred while sending the email!');
          return;
        }
        return handleSendEmail(
          data.replace(/(<([^>]+)>)/gi, ''),
          () => console.log(''),
          to,
          respondTo,
        );
      },
    },
  ];

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: isEditMode ? '-17px -24px' : 0,
      }}
    >
      <SocialEditor
        editorRef={editorRef}
        mode={NoteEditorModes.ADD}
        saving={saving}
        value={''}
        manager={manager}
        state={state}
        setState={setState}
        // handleUploadClick={}
        items={
          editorModeState.mode === EditorMode.Email
            ? submitEmailButtonOptions
            : submitButtonOptions
        }
      />
    </div>
  );
};

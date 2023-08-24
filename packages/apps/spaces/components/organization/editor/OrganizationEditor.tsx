import React, { FC, useCallback, useRef } from 'react';
import { editorEmail, editorMode } from '../../../state';
import { useRecoilState, useRecoilValue } from 'recoil';
import {
  extraAttributes,
  SocialEditor,
} from '@ui/form/RichTextEditor/SocialEditor';
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
import { NoteEditorModes } from './types';

interface Props {
  mode: NoteEditorModes;
  organizationId: string;
}
export const OrganizationEditor: FC<Props> = ({
  mode = NoteEditorModes.ADD,
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
  const {
    handleSubmit: handleSendEmail,
    to,
    respondTo,
  } = useRecoilValue(editorEmail);

  const editorRef = useRef<any | null>(null);

  const isEditMode = mode === NoteEditorModes.EDIT;
  const handleResetEditor = () => {
    const context = getContext();
    if (context) {
      context.commands.resetContent();
    }
  };

  const handleSendEmailResponse = () => {
    const data = prosemirrorNodeToHtml(state.doc);
    if (!handleSendEmail) {
      toast.error('Client error occurred while sending the email!');
      return;
    }
    return handleSendEmail(
      data.replace(/(<([^>]+)>)/gi, ''),
      () => handleResetEditor(),
      to,
      respondTo,
    );
  };

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
        value={''}
        manager={manager}
        state={state}
        setState={setState}
        onSubmit={handleSendEmailResponse}
        submitButtonLabel={'Send'}
      />
    </div>
  );
};

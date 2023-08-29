import React from 'react';
import { Remirror, useRemirror } from '@remirror/react';
import { AnyExtension, RemirrorManager } from 'remirror';

interface RichTextPreviewProps {
  htmlContent: string;
  extensions: RemirrorManager<any> | (() => AnyExtension[]) | undefined;
}

export const RichTextPreview: React.FC<RichTextPreviewProps> = ({
  htmlContent,
  extensions,
}) => {
  const { state, manager } = useRemirror({
    extensions,
    content: htmlContent,
    stringHandler: 'html',
  });

  return <Remirror manager={manager} initialContent={state} editable={false} />;
};

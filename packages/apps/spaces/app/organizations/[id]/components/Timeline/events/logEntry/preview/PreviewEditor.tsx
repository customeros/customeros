import React, { useEffect } from 'react';
import { Text } from '@ui/typography/Text';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { TagSuggestor } from '@ui/form/RichTextEditor/TagSuggestor';
import { useRemirror } from '@remirror/react';
import { logEntryEditorExtensions } from '@organization/components/Timeline/TimelineActions/context/extensions';
import { RichEditorBlurHandler } from '@ui/form/RichTextEditor/components/RichEditorBlurHandler';
import { htmlToProsemirrorNode } from 'remirror';

export const PreviewEditor: React.FC<{
  isAuthor: boolean;
  formId: string;
  initialContent: string;
  tags?: Array<{ label: string; value: string }>;
}> = ({ isAuthor, formId, initialContent, tags }) => {
  const remirrorProps = useRemirror({
    extensions: logEntryEditorExtensions,
  });
  useEffect(() => {
    const prosemirrorNodeValue = htmlToProsemirrorNode({
      schema: remirrorProps.state.schema,
      content: `${initialContent}`,
    });
    remirrorProps.getContext()?.setContent(prosemirrorNodeValue);
  }, [initialContent]);

  return (
    <>
      {!isAuthor && (
        <HtmlContentRenderer
          fontSize='sm'
          noOfLines={undefined}
          htmlContent={`${initialContent}`}
        />
      )}

      {isAuthor && (
        <Text sx={{ '& .remirror-theme': { fontSize: 'sm' } }}>
          <RichTextEditor
            {...remirrorProps}
            placeholder='Log a conversation you had with a customer'
            formId={formId}
            name='content'
            showToolbar={false}
          >
            <>
              <TagSuggestor
                tags={tags?.map((e: { label: string; value: string }) => ({
                  label: e.label,
                  id: e.value,
                }))}
              />
              <RichEditorBlurHandler formId={formId} name='content' />
            </>
          </RichTextEditor>
        </Text>
      )}
    </>
  );
};

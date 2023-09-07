import React, {
  FC,
  forwardRef,
  PropsWithChildren,
  useCallback,
  useImperativeHandle,
} from 'react';
import {
  Remirror,
  ThemeProvider,
  Toolbar,
  useExtension,
} from '@remirror/react';
import { useField } from 'react-inverted-form';
import { prosemirrorNodeToHtml } from 'remirror';
import {
  BasicEditorExtentions,
  RemirrorProps,
} from '@ui/form/RichTextEditor/types';
import { MentionAtomExtension } from 'remirror/extensions';

export const TagArea: FC<
  {
    name: string;
    formId: string;
    showToolbar: boolean;
    submit: any;
  } & RemirrorProps<BasicEditorExtentions> &
    PropsWithChildren
> = forwardRef(
  (
    {
      children,
      name,
      showToolbar,
      formId,
      manager,
      getContext,
      state,
      setState,
      submit,
    },
    ref,
  ) => {
    const { getInputProps } = useField(name, formId);
    const { onChange, value } = getInputProps();
    useImperativeHandle(ref, () => getContext(), [getContext]);

    return (
      <ThemeProvider>
        <Remirror
          manager={manager}
          placeholder='Log conversation you had with a customer'
          onChange={(parameter) => {
            const nextState = parameter.state;
            const htmlValue = prosemirrorNodeToHtml(nextState?.doc);

            // first update is happening before form store is initialized this change prevents error
            if (value !== undefined) {
              onChange(htmlValue);
            }
            setState(nextState);
          }}
          initialContent={state}
          autoRender='end'
        >
          {showToolbar && (
            <Toolbar
              height={'var(--chakra-sizes-8)'}
              style={{ overflowX: 'visible' }}
            >
              {children}
            </Toolbar>
          )}
          {submit && submit}
        </Remirror>
      </ThemeProvider>
    );
  },
);

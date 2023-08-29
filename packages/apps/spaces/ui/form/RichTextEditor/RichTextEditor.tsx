import React, {
  FC,
  forwardRef,
  PropsWithChildren,
  useImperativeHandle,
} from 'react';
import { Remirror, ThemeProvider, Toolbar } from '@remirror/react';
import { useField } from 'react-inverted-form';
import { prosemirrorNodeToHtml } from 'remirror';
import {
  BasicEditorExtentions,
  RemirrorProps,
} from '@ui/form/RichTextEditor/types';

export const RichTextEditor: FC<
  {
    name: string;
    formId: string;
  } & RemirrorProps<BasicEditorExtentions> &
    PropsWithChildren
> = forwardRef(
  ({ children, name, formId, manager, getContext, state, setState }, ref) => {
    const { getInputProps } = useField(name, formId);
    const { onChange, value } = getInputProps();
    useImperativeHandle(ref, () => getContext(), [getContext]);

    return (
      <ThemeProvider>
        <Remirror
          manager={manager}
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
          <Toolbar
            height={'var(--chakra-sizes-8)'}
            style={{ overflowX: 'visible' }}
          >
            {children}
          </Toolbar>
        </Remirror>
      </ThemeProvider>
    );
  },
);

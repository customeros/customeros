import { useField } from 'react-inverted-form';
import {
  FC,
  useRef,
  useEffect,
  forwardRef,
  useCallback,
  PropsWithChildren,
  useImperativeHandle,
} from 'react';

import { prosemirrorNodeToHtml } from 'remirror';
import { ClickHandlerState } from '@remirror/extension-events';
import {
  Remirror,
  useCommands,
  ThemeProvider,
  useEditorEvent,
} from '@remirror/react';

import { RemirrorProps } from '@ui/form/RichTextEditor/types';
import { FloatingLinkToolbar } from '@ui/form/RichTextEditor/floatingMenu/FloatingLinkMenu';

const hooks = [
  () => {
    const { selectText } = useCommands();

    const clickHandler = useCallback(
      (e: MouseEvent, clickHandlerState: ClickHandlerState) => {
        // @ts-expect-error nodeName exists but type is not compatible
        if (e?.target?.nodeName === 'A') {
          selectText(clickHandlerState.markRanges[0]);
        }
      },
      [selectText],
    );
    useEditorEvent('click', clickHandler);
  },
];
export const RichTextEditor: FC<
  {
    name: string;
    formId: string;
    placeholder?: string;
    showToolbar: boolean;
    // exhaustively typing this is not really necessary for us at the moment
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } & RemirrorProps<any> &
    PropsWithChildren
> = forwardRef(
  (
    {
      children,
      name,
      formId,
      manager,
      getContext,
      state,
      setState,
      placeholder = '',
    },
    ref,
  ) => {
    const didMountRef = useRef(false);
    const { getInputProps } = useField(name, formId);
    const { onChange, value } = getInputProps();
    useImperativeHandle(ref, () => getContext(), [getContext]);

    // TODO: remove this when react-inverted-form will prevent handler calls before form is initialized completely
    useEffect(() => {
      if (didMountRef.current) {
        return;
      }
      didMountRef.current = true;
    }, []);

    return (
      <ThemeProvider>
        <Remirror
          manager={manager}
          placeholder={placeholder}
          onChange={(parameter) => {
            const nextState = parameter.state;
            const htmlValue = prosemirrorNodeToHtml(nextState?.doc);
            // first update is happening before form store is initialized this change prevents error
            if (value !== undefined && didMountRef.current) {
              onChange?.(htmlValue);
            }
            setState(nextState);
          }}
          initialContent={state}
          autoRender='end'
          hooks={hooks}
        >
          <FloatingLinkToolbar />
          {children}
        </Remirror>
      </ThemeProvider>
    );
  },
);

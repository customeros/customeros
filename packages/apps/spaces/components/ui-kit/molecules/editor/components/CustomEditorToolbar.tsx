import type { ChangeEvent, HTMLProps, KeyboardEvent } from 'react';
import React, { useCallback, useEffect, useRef } from 'react';

import {
  CommandButton,
  FloatingToolbar,
  FloatingWrapper,
  ToggleBoldButton,
  ToggleItalicButton,
  ToggleStrikeButton,
  ToggleUnderlineButton,
  useActive,
  useCurrentSelection,
} from '@remirror/react';
import { useFloatingLinkState } from '../hooks';

const DelayAutoFocusInput = ({
  autoFocus,
  ...rest
}: HTMLProps<HTMLInputElement>) => {
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!autoFocus) {
      return;
    }

    const frame = window.requestAnimationFrame(() => {
      inputRef.current?.focus();
    });

    return () => {
      window.cancelAnimationFrame(frame);
    };
  }, [autoFocus]);

  return <input ref={inputRef} {...rest} />;
};

export const CustomEditorToolbar = ({ editable = false }) => {
  const {
    isEditing,
    linkPositioner,
    clickEdit,
    onRemove,
    submitHref,
    href,
    setHref,
    cancelHref,
  } = useFloatingLinkState();
  const active = useActive();
  const activeLink = active.link();
  const { empty } = useCurrentSelection();

  const handleClickEdit = useCallback(() => {
    clickEdit();
  }, [clickEdit]);

  const linkEditButtons = activeLink ? (
    <>
      <CommandButton
        commandName='updateLink'
        onSelect={handleClickEdit}
        icon='pencilLine'
        enabled
      />
      <CommandButton
        commandName='removeLink'
        onSelect={onRemove}
        icon='linkUnlink'
        enabled
      />
    </>
  ) : (
    <CommandButton
      commandName='updateLink'
      onSelect={handleClickEdit}
      icon='link'
      enabled
    />
  );

  return (
    <>
      {!isEditing && editable && (
        <FloatingToolbar>
          {linkEditButtons}
          <ToggleBoldButton />
          <ToggleItalicButton />
          <ToggleUnderlineButton />
          <ToggleStrikeButton />
        </FloatingToolbar>
      )}
      {!isEditing && empty && editable && (
        <FloatingToolbar>{linkEditButtons}</FloatingToolbar>
      )}

      <FloatingWrapper
        positioner='always'
        placement='bottom'
        enabled={isEditing}
        renderOutsideEditor
      >
        <DelayAutoFocusInput
          style={{ zIndex: 20 }}
          autoFocus
          placeholder='Enter link...'
          onChange={(event: ChangeEvent<HTMLInputElement>) =>
            setHref(event.target.value)
          }
          value={href}
          onKeyDown={(event: KeyboardEvent<HTMLInputElement>) => {
            const { code } = event;

            if (code === 'Enter') {
              submitHref();
            }

            if (code === 'Escape') {
              cancelHref();
            }
          }}
        />
      </FloatingWrapper>
    </>
  );
};

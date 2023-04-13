import React, { useEffect, useState, useRef } from 'react';
import styles from './editable-content-input.module.scss';
import classNames from 'classnames';
import { useDebouncedCallback } from 'use-debounce';
import { ChevronDown, ChevronUp } from '../icons';
import { IconButton } from '../icon-button';

export const DebouncedTextArea = ({
  placeholder = '',
  inputSize = 'md',
  debounceTimeout = 300,
  value = '',
  onChange,
  isEditMode,
  ...props
}: any) => {
  const [inner, setInner] = useState(value);
  const [expanded, setExpanded] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const debounced = useDebouncedCallback(
    // function
    (value) => {
      onChange(value);
    },
    // delay in ms
    debounceTimeout,
  );

  useEffect(() => {
    if (!isEditMode) {
      debounced.flush();
    }
  }, [isEditMode]);

  useEffect(() => {
    if (value && !isEditMode) {
      setInner(value);
    }
    if (!isEditMode) {
      setExpanded(false);
    }
  }, [isEditMode, value]);

  useEffect(() => {
    if (isEditMode && textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  }, [isEditMode, inner]);

  useEffect(() => {
    if (!isEditMode && expanded && textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  }, [isEditMode, expanded]);

  console.log(
    '🏷️ ----- :test ',
    !isEditMode &&
      !expanded &&
      textareaRef?.current &&
      textareaRef.current.scrollHeight > 320,
  );
  return (
    <div className={styles.textAreaContainer}>
      <textarea
        {...props}
        ref={textareaRef}
        value={inner}
        className={classNames(styles.contentEditable, styles.textArea, {
          [styles?.[inputSize]]: inputSize,
          [styles.editable]: isEditMode,
          [styles.readMode]: !isEditMode && expanded,
          [styles.collapsed]: !isEditMode && !expanded,
          [styles.readModeAllowScroll]:
            !isEditMode &&
            expanded &&
            textareaRef?.current &&
            textareaRef.current.scrollHeight > 320,
        })}
        readOnly={!isEditMode}
        onChange={(event) => {
          setInner(event.target.value);
          debounced(event.target.value);
          if (isEditMode && textareaRef.current) {
            textareaRef.current.style.height = 'auto';
            textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
          }
        }}
        placeholder={placeholder}
        onBlur={debounced.flush}
      />
      <div
        className={classNames({
          [styles.blurVisible]: !isEditMode && !expanded,
        })}
      />
      {!isEditMode && (
        <IconButton
          className={styles.collapseExpandButton}
          mode='text'
          size='xxxs'
          onClick={() => setExpanded(!expanded)}
          icon={expanded ? <ChevronUp /> : <ChevronDown />}
        />
      )}
    </div>
  );
};

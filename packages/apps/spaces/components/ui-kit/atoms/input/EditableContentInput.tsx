import React, { useEffect, useRef, useState } from 'react';
import styles from './editable-content-input.module.scss';
import classNames from 'classnames';
import { useDebouncedCallback } from 'use-debounce';
export const EditableContentInput = ({
  placeholder = '',
  inputSize = 'md',
  debounceTimeout = 300,
  value = '',
  onChange,
  isEditMode,
  ...rest
}: any) => {
  const inputRef = useRef<Node>(null);
  const [inner, setInner] = useState(value);

  const debounced = useDebouncedCallback(
    // function
    (value) => {
      onChange(value);
    },
    // delay in ms
    debounceTimeout,
  );

  useEffect(() => {
    return () => {
      debounced.flush();
    };
  }, []);

  return (
    <input
      {...rest}
      value={inner}
      size={inner.length || placeholder.length}
      ref={inputRef}
      className={classNames(styles.contentEditable, {
        [styles?.[inputSize]]: inputSize,
        [styles.editable]: isEditMode,
      })}
      disabled={!isEditMode}
      style={{ textAlign: inner.length ? 'center' : 'initial' }}
      onChange={(event) => {
        setInner(event.target.value);
        debounced(event.target.value);
      }}
      placeholder={placeholder}
    />
  );
};

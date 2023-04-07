import React, { useEffect, useState } from 'react';
import styles from './editable-content-input.module.scss';
import classNames from 'classnames';
import { useDebouncedCallback } from 'use-debounce';
export const DebouncedTextArea = ({
  placeholder = '',
  inputSize = 'md',
  debounceTimeout = 300,
  value = '',
  onChange,
  isEditMode,
  ...rest
}: any) => {
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

  if (!isEditMode) {
    return (
      <div
        {...rest}
        className={classNames(styles.contentEditable, styles.textArea, {
          [styles?.[inputSize]]: inputSize,
          [styles.editable]: isEditMode,
        })}
      >
        {value}
      </div>
    );
  }

  return (
    <>
      <textarea
        {...rest}
        value={inner}
        className={classNames(styles.contentEditable, styles.textArea, {
          [styles?.[inputSize]]: inputSize,
          [styles.editable]: isEditMode,
        })}
        disabled={!isEditMode}
        onChange={(event) => {
          setInner(event.target.value);
          debounced(event.target.value);
        }}
        placeholder={placeholder}
      />
    </>
  );
};

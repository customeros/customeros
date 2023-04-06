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
  const inputRef = useRef<HTMLSpanElement>(null);
  const [inner, setInner] = useState(value);
  const [width, setWidth] = useState<number>();

  const debounced = useDebouncedCallback(
    // function
    (value) => {
      onChange(value);
    },
    // delay in ms
    debounceTimeout,
  );
  useEffect(() => {
    setWidth((inputRef?.current?.scrollWidth || 0) + 2);
  }, [inner]);

  useEffect(() => {
    return () => {
      debounced.flush();
    };
  }, []);

  return (
    <>
      <input
        {...rest}
        value={inner}
        size={inner.length || placeholder.length}
        className={classNames(styles.contentEditable, {
          [styles?.[inputSize]]: inputSize,
          [styles.editable]: isEditMode,
        })}
        style={{ width: `${width}px` }}
        disabled={!isEditMode}
        onChange={(event) => {
          setInner(event.target.value);
          debounced(event.target.value);
        }}
        placeholder={placeholder}
      />
      <span
        ref={inputRef}
        className={classNames(styles.contentEditable, {
          [styles?.[inputSize]]: inputSize,
          [styles.editable]: isEditMode,
        })}
        style={{ visibility: 'hidden', position: 'absolute' }}
      >
        {inner || placeholder}
      </span>
    </>
  );
};

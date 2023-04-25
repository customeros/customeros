import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
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
  label,
  id,
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
  useLayoutEffect(() => {
    setWidth((inputRef?.current?.scrollWidth || 0) + 2);
  }, [inner, isEditMode]);

  useEffect(() => {
    return () => {
      debounced.flush();
    };
  }, []);

  return (
    <>
      <label
        htmlFor={`editable-content-input-${id}`}
        className={styles.invisibleLabel}
      >
        {label}
      </label>
      <input
        {...rest}
        id={`editable-content-input-${id}`}
        value={inner}
        size={inner.length || placeholder.length}
        className={classNames(styles.contentEditable, {
          [styles?.[inputSize]]: inputSize,
          [styles.editable]: isEditMode,
          [rest.className]: rest.className,
        })}
        disabled={!isEditMode}
        style={{ width: `${width}px` }}
        onChange={(event) => {
          setInner(event.target.value);
          debounced(event.target.value);
        }}
        placeholder={placeholder}
        onBlur={() => debounced.flush()}
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

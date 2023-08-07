import React, { CSSProperties, FC } from 'react';
import { useSelect } from '../useSelect';
import classNames from 'classnames';
import { InlineLoader } from '@ui/presentation/inline-loader';
import styles from './select.module.scss';

export const SelectInput: FC<{
  saving?: boolean;
  readOnly?: boolean;
  placeholder: string;
  customStyles?: CSSProperties | undefined;
}> = ({ saving, placeholder, customStyles, readOnly }) => {
  const { state, getInputProps, autofillValue } = useSelect();
  return (
    <>
      <span
        role='textbox'
        placeholder={placeholder}
        contentEditable={state.isEditing && !readOnly}
        className={classNames(styles.dropdownInput, {
          [styles.dropdownInputEditable]: state.isEditing && !readOnly,
        })}
        style={customStyles}
        {...getInputProps()}
      />
      <span
        className={classNames(styles.autofill, {
          [styles.autofillIndentation]:
            autofillValue?.charAt(0) === ' ' &&
            autofillValue?.charAt(1) !== ' ',
        })}
      >
        {autofillValue}
      </span>
      {saving && <InlineLoader />}
    </>
  );
};

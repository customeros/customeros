import React from 'react';
import styles from './simple-validation-indicator.module.scss';
import CheckCircleFilled from '@spaces/atoms/icons/CheckCircleFilled';
import { InlineLoader } from '@spaces/atoms/inline-loader';

interface Props {
  showValidationMessage: boolean;
  isLoading: boolean;
  isEditMode: boolean;
  errorMessages?: Array<string>;
}

export const SimpleValidationIndicator = ({
  showValidationMessage,
  isLoading,
  isEditMode,
  errorMessages = [],
}: Props) => {
  if (isLoading) {
    return <InlineLoader label='Validating' color='#DB9E00' />;
  }

  if (!errorMessages.length && isEditMode) {
    return (
      <div className={styles.validEntry}>
        <CheckCircleFilled color='#BED9B1' height={8} width={8} />
      </div>
    );
  }

  if (!errorMessages.length && !isEditMode) {
    return null;
  }

  return (
    <div className={styles.validationSignal}>
      {showValidationMessage && (
        <div className={styles.validationMessage}>
          {errorMessages.map((data) => (
            <div key={data.split(' ').join('-')}>{data}</div>
          ))}
        </div>
      )}
    </div>
  );
};

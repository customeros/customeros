import React, { ButtonHTMLAttributes, FC, ReactNode } from 'react';
import { Dialog } from 'primereact/dialog';
import { Button } from '../button';
import styles from './delete-confirmation-dialog.module.scss';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  deleteConfirmationModalVisible: boolean;
  setDeleteConfirmationModalVisible: (newState: boolean) => void;
  deleteAction: () => void;
  header?: string;
  confirmationButtonLabel?: string;
  explanationText?: string | ReactNode;
}

export const DeleteConfirmationDialog: FC<Props> = ({
  setDeleteConfirmationModalVisible,
  deleteConfirmationModalVisible,
  deleteAction,
  header,
  confirmationButtonLabel,
  explanationText,
}) => {
  return (
    <>
      <Dialog
        header={header || 'Confirm Delete'}
        draggable={false}
        className={styles.dialog}
        visible={deleteConfirmationModalVisible}
        footer={
          <div className={styles.dialogFooter}>
            <Button
              mode='secondary'
              onClick={() => setDeleteConfirmationModalVisible(false)}
              className='p-button-text'
            >
              Cancel
            </Button>
            <Button
              style={{ marginRight: '0' }}
              mode='danger'
              onClick={deleteAction}
              autoFocus
            >
              {confirmationButtonLabel || 'Delete'}
            </Button>
          </div>
        }
        onHide={() => setDeleteConfirmationModalVisible(false)}
      >
        <p className={styles.dialogExplanation}>
          {explanationText ||
            'Are you sure you want to delete this item? This action cannot be undone.'}
        </p>
      </Dialog>
    </>
  );
};

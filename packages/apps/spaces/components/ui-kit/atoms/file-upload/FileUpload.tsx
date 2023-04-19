import React, { useState } from 'react';
import styles from './file-upload.module.scss';
import classNames from 'classnames';
import { Paperclip } from '../icons';

export const FileUpload = ({ onFileUpload, className, children }: any) => {
  const [isDraggingOver, setIsDraggingOver] = useState(false);

  const handleDrop = (ev: any) => {
    console.log('File(s) dropped');

    // Prevent default behavior (Prevent file from being opened)
    ev.preventDefault();

    if (ev.dataTransfer.items) {
      // Use DataTransferItemList interface to access the file(s)
      [...ev.dataTransfer.items].forEach((item, i) => {
        // If dropped items aren't files, reject them
        if (item.kind === 'file') {
          const file = item.getAsFile();
          console.log(`… file[${i}].name = ${file.name}`);
        }
      });
    } else {
      // Use DataTransfer interface to access the file(s)
      [...ev.dataTransfer.files].forEach((file, i) => {
        console.log(`… file[${i}].name = ${file.name}`);
      });
    }
  };

  return (
    <section
      className={classNames(styles.fileUploadContainer, {
        [className]: !!className,
        [styles.dragOver]: isDraggingOver,
      })}
      onDrop={handleDrop}
      onDragOver={() => setIsDraggingOver(true)}
      onDragLeave={() => setIsDraggingOver(false)}
    >
      <h3 className={styles.attachmentTitle}>Attach a file</h3>
      <Paperclip />
    </section>
  );
};

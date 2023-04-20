import React, { useState } from 'react';
import styles from './file-upload.module.scss';
import classNames from 'classnames';
import { Paperclip } from '../icons';
import axios from 'axios';
import { toast } from 'react-toastify';
import Image from 'next/image';
import { IconButton } from '../icon-button';
import { uuid4 } from '@sentry/utils';

export const FileUpload = ({
  files,
  onBeginFileUpload,
  onFileUpload,
  onFileUploadError,
  onFileRemove,
}: any) => {
  const uploadInputRef = React.useRef<HTMLInputElement>(null);
  const [isDraggingOver, setIsDraggingOver] = useState(false);

  const handleDrag = function (e: any) {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setIsDraggingOver(true);
    } else if (e.type === 'dragleave') {
      setIsDraggingOver(false);
    }
  };

  const handleDrop = (ev: any) => {
    // Prevent default behavior (Prevent file from being opened)
    ev.preventDefault();
    ev.stopPropagation();

    setIsDraggingOver(false);

    if (ev.dataTransfer.items) {
      // Use DataTransferItemList interface to access the file(s)
      [...ev.dataTransfer.items].forEach((item, i) => {
        // If dropped items aren't files, reject them
        if (item.kind === 'file') {
          const file = item.getAsFile();
          handleInputFileChange({ target: { files: [file] } });
        }
      });
    } else {
      // Use DataTransfer interface to access the file(s)
      [...ev.dataTransfer.files].forEach((file, i) => {
        handleInputFileChange({ target: { files: [file] } });
      });
    }
  };

  const handleInputFileChange = (e: any) => {
    if (!e?.target?.files) {
      return;
    }

    const fileKey = uuid4();
    onBeginFileUpload(fileKey);

    const formData = new FormData();
    formData.append('file', e.target.files[0]);

    const clearFileInput = () => {
      uploadInputRef && uploadInputRef.current
        ? (uploadInputRef.current.value = '')
        : '';
    };

    axios
      .post(`/fs/file`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      .then((r: any) => {
        onFileUpload({ ...r.data, key: fileKey });

        clearFileInput();

        return r.data;
      })
      .catch((e) => {
        clearFileInput();
        onFileUploadError(fileKey);
        toast.error(
          'Oops! We could add this file. Check if file type is supported and can try again or contact our support team',
        );
      });
  };

  const fileTemplate = (file: any, fileType: string) => {
    console.log(fileType);
    return (
      <div key={file.id + '_' + file.key} className={styles.fileContainer}>
        <div
          className={styles.removeFile}
          onClick={() => onFileRemove(file.id)}
        >
          <IconButton
            aria-describedby='message-icon-label'
            mode='secondary'
            className={styles.icon}
            onClick={() => null}
            icon={
              <Image alt={''} src='/icons/ban.svg' width={20} height={20} />
            }
          />
        </div>
        <div className={styles.preview}>
          {(fileType == undefined ||
            (fileType !== 'png' &&
              fileType !== 'jpg' &&
              fileType !== 'jpeg')) && (
            <Image alt={''} src='/icons/file.svg' width={75} height={75} />
          )}

          {(fileType == 'png' || fileType == 'jpg' || fileType == 'jpeg') && (
            <Image alt={''} src='/icons/image.svg' width={75} height={75} />
          )}
        </div>
        <div className={styles.text}>
          {!file.uploaded && <>In progress</>}

          {file.uploaded && file.name}
        </div>
      </div>
    );
  };
  return (
    <section
      className={classNames(styles.fileUploadContainer, {
        [styles.dragOver]: isDraggingOver,
      })}
      onDragEnter={handleDrag}
      onDragLeave={handleDrag}
      onDragOver={handleDrag}
      onDrop={handleDrop}
    >
      <div className={styles.files}>
        {files?.length > 0 &&
          files.map((file: any) => {
            return fileTemplate(file, file.extension);
          })}
      </div>
      <div
        className={styles.attachFile}
        onClick={() => uploadInputRef?.current?.click()}
      >
        <h3 className={styles.attachFileText}>Attach a file</h3>
        <Paperclip />

        <input
          style={{ display: 'none' }}
          ref={uploadInputRef}
          type='file'
          onChange={handleInputFileChange}
        />
      </div>
    </section>
  );
};

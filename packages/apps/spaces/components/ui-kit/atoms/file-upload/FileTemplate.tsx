import React, { FC } from 'react';
import styles from './file-upload.module.scss';

import Image from 'next/image';
import { DeleteIconButton, IconButton } from '../icon-button';
import { Skeleton } from '../skeleton';

interface FileTemplateProps {
  file: any;
  fileType: string;
  onFileRemove: (id: string) => void;
}
export const FileTemplate: FC<FileTemplateProps> = ({
  file,
  fileType,
  onFileRemove,
}) => {
  return (
    <div key={file.id + '_' + file.key} className={styles.fileContainer}>
      <div className={styles.removeFile}>
        <DeleteIconButton onDelete={() => onFileRemove(file.id)} />
      </div>
      <div className={styles.preview}>
        {(fileType == undefined ||
          (fileType !== 'png' &&
            fileType !== 'jpg' &&
            fileType !== 'jpeg')) && (
          <Image alt={''} src='/icons/file.svg' width={40} height={40} />
        )}

        {(fileType == 'png' || fileType == 'jpg' || fileType == 'jpeg') && (
          <Image alt={''} src='/icons/image.svg' width={40} height={40} />
        )}
      </div>
      <div className={styles.text}>
        {!file.uploaded && <Skeleton height={'5px'} />}

        {file.uploaded && file.name}
      </div>
    </div>
  );
};

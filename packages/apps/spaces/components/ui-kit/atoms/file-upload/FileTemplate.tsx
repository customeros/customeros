import React, { FC } from 'react';
import styles from './file-upload.module.scss';

import { DeleteIconButton } from '../icon-button';
import { Skeleton } from '../skeleton';
import { File, Image } from '../icons';

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
    <div
      key={file.id + '_' + file.key}
      className={styles.fileContainer}
      title={`${file.name}.${fileType}`}
    >
      <div className={styles.removeFile}>
        <DeleteIconButton onDelete={() => onFileRemove(file.id)} />
      </div>
      <div className={styles.preview}>
        {(fileType == undefined ||
          (fileType !== 'png' &&
            fileType !== 'jpg' &&
            fileType !== 'jpeg')) && (
          <File
            height={24}
            width={24}
            aria-label={file.uploaded ? '' : 'Uploading file'}
          />
        )}

        {(fileType == 'png' || fileType == 'jpg' || fileType == 'jpeg') && (
          // eslint-disable-next-line jsx-a11y/alt-text
          <Image
            width={24}
            height={24}
            aria-label={file.uploaded ? '' : 'Uploading file'}
          />
        )}
      </div>
      <div className={styles.text}>
        {!file.uploaded && !file.name && <Skeleton height={'5px'} />}

        {file.name}
      </div>
    </div>
  );
};

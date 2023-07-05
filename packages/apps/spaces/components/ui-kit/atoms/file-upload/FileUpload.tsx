import React, { useState } from 'react';
import styles from './file-upload.module.scss';
import classNames from 'classnames';
import { CloudUpload } from '../icons';
import axios from 'axios';
import { toast } from 'react-toastify';
import { uuid4 } from '@sentry/utils';
import { Flex, Text } from '@chakra-ui/react';

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

  return (
    <Flex
      className={classNames(styles.drag, {
        [styles.dragOver]: isDraggingOver,
      })}
      alignItems='center'
      justifyContent='center'
      width='100%'
      direction='column'
      onDragEnter={handleDrag}
      onDragLeave={handleDrag}
      onDragOver={handleDrag}
      onDrop={handleDrop}
      padding='md'
      background='#fff'
      borderRadius='xl'
      border='1px solid #EAECF0'
      mt='2'
    >
      <div className={styles.iconUpload}>
        <CloudUpload color='#000' height='20px' width='20px' />
      </div>
      <div className={styles.attachFile}>
        <Text
          color='#6941C6'
          size='sm'
          fontWeight={600}
          as='button'
          mr={1}
          onClick={() => uploadInputRef?.current?.click()}
        >
          Click to upload
        </Text>

        <Text color='gray.600' size='sm' as='span'>
          or drag and drop
        </Text>

        <Text color='gray.600' fontSize='14px' size='xs' textAlign='center'>
          Max. 20MB in size
        </Text>

        <input
          style={{ display: 'none' }}
          ref={uploadInputRef}
          type='file'
          onChange={handleInputFileChange}
        />
      </div>
    </Flex>
  );
};

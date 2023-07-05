import React, { FC } from 'react';
import styles from './file-upload.module.scss';

import { DeleteIconButton } from '../icon-button';
import { Skeleton } from '../skeleton';
import { File, Image } from '../icons';
import { Badge, Spinner } from '@chakra-ui/react';
import { IconButton } from '@ui/form/IconButton';
import TimesOutlined from '@spaces/atoms/icons/TimesOutlined';
import { Loader } from '@spaces/atoms/loader';
import Paperclip from '@spaces/atoms/icons/Paperclip';

interface FileTemplateProps {
  file: any;
  fileType: string;
  onFileRemove: (id: string) => void;
}
export const FileTemplateUpload: FC<FileTemplateProps> = ({
  file,
  fileType,
  onFileRemove,
}) => {
  return (
    <>
      {(fileType == undefined ||
        (fileType !== 'png' && fileType !== 'jpg' && fileType !== 'jpeg')) && (
        <Badge
          boxShadow='none'
          variant='outline'
          borderRadius='xl'
          borderWidth='1px'
          fontWeight={300}
          textTransform='initial'
          alignItems='center'
          display='inline-flex'
          background='gray.50'
          px={2}
          mr={1}
        >
          {file.name}
          {fileType}

          {file.uploaded && (
            <IconButton
              variant='ghost'
              aria-label='Close preview'
              color='gray.500'
              borderRadius={30}
              padding={0}
              height='12px'
              width='12px'
              minWidth='12px'
              minHeight='12px'
              ml={1}
              icon={<TimesOutlined color='#98A2B3' height='12px' />}
              onClick={() => onFileRemove(file.id)}
            />
          )}

          {!file.uploaded && (
            <Spinner size='xs' color='#7F56D9' ml={1} speed='1s' />
          )}
        </Badge>
      )}
    </>
  );
};

export const ReadAttachmentBadge: FC<FileTemplateProps> = ({
  file,
  fileType,
  onFileRemove,
}) => {
  return (
    <Badge
      boxShadow='none'
      variant='outline'
      borderRadius='xl'
      borderWidth='1px'
      fontWeight={300}
      textTransform='initial'
      alignItems='center'
      display='inline-flex'
      background='gray.50'
      px={2}
      mr={1}
    >
      <Paperclip height={12} width={12} />
      {file.name}
      {fileType}
    </Badge>
  );
};

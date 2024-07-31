import { forwardRef } from 'react';

import { useFileUploader, FileUploaderProps } from './useFileUploader';

interface FileUploadTriggerProps extends FileUploaderProps {
  name?: string;
  accept?: string;
  className?: string;
  children?: React.ReactNode;
}

export const FileUploadTrigger = forwardRef<
  HTMLLabelElement,
  FileUploadTriggerProps
>(
  (
    {
      accept,
      onError,
      onChange,
      onLoading,
      onSuccess,
      onLoadEnd,
      onProgress,
      apiBaseUrl,
      onLoadStart,
      endpointOptions,
      name = 'fileUploadInput',
      className,
      ...props
    },
    ref,
  ) => {
    const { inputRef, handleOnChange } = useFileUploader({
      onError,
      onChange,
      onLoadEnd,
      onLoading,
      onSuccess,
      onProgress,
      apiBaseUrl,
      onLoadStart,
      endpointOptions,
    });

    return (
      <>
        <label ref={ref} htmlFor={name} className={className} {...props} />

        <input
          id={name}
          type='file'
          ref={inputRef}
          accept={accept}
          className='hidden'
          onChange={handleOnChange}
        />
      </>
    );
  },
);

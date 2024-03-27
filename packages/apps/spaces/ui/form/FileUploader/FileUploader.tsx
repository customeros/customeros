import { useRef, forwardRef, ChangeEventHandler } from 'react';

export interface EndpointOptions {
  uploadUrl: string;
  fileKeyName: string;
  downloadUrl?: string;
  payload?: Record<string, string>;
}

interface FileUploaderProps {
  name?: string;
  accept?: string;
  apiBaseUrl: string;
  className?: string;
  children?: React.ReactNode;
  onLoadStart?(refId: number): void;
  endpointOptions?: EndpointOptions;
  onChange?(file: File, refId: number): void;
  onLoadEnd?(refId: number, e: unknown): void;
  onError?(refId: number, error: string): void;
  onLoading?(refId: number, value: boolean): void;
  onProgress?(refId: number, value: number): void;
  onSuccess?(refId: number, response: unknown): void;
}

export const FileUploader = forwardRef<HTMLLabelElement, FileUploaderProps>(
  (
    {
      name = 'fileUnploadInput',
      apiBaseUrl,
      endpointOptions,
      accept,
      onChange,
      onLoading,
      onProgress,
      onLoadStart,
      onLoadEnd,
      onSuccess,
      onError,
      className,
      ...props
    },
    ref,
  ) => {
    const inputRef = useRef<HTMLInputElement>(null);

    const handleOnChange: ChangeEventHandler<HTMLInputElement> = async () => {
      if (!inputRef) return;

      const currentFiles = inputRef.current?.files;

      if (!currentFiles) return;

      for (const file of currentFiles) {
        const refId = Math.random();
        onChange?.(file, refId);

        const xhr = new XMLHttpRequest();

        xhr.upload.onloadstart = () => {
          onLoading?.(refId, true);
          onLoadStart?.(refId);
        };

        xhr.upload.onprogress = (e) => {
          const percentage = Math.floor((e.loaded * 100) / e.total);
          onProgress?.(refId, percentage);
        };

        xhr.upload.onerror = () => {
          onLoading?.(refId, false);
          onError?.(refId, 'Could not upload file.');
        };

        xhr.upload.onabort = () => {
          onLoading?.(refId, false);
        };

        xhr.upload.ontimeout = () => {
          onLoading?.(refId, false);
        };

        xhr.onloadend = (e) => {
          if (xhr.status >= 200 && xhr.status < 300) {
            onLoadEnd?.(refId, e);
          } else {
            onError?.(refId, 'Could not upload file.');
          }

          onLoading?.(refId, false);
        };

        xhr.onreadystatechange = () => {
          if (xhr.readyState === 4 && xhr.status === 200) {
            const data = JSON.parse(xhr.responseText);
            onSuccess?.(refId, data);
          } else if (xhr.readyState === 4) {
            onError?.(refId, 'Could not fetch data.');
          }
        };

        if (!endpointOptions) {
          return;
        }

        const formData = new FormData();

        if (endpointOptions) {
          xhr.open('POST', `${apiBaseUrl}${endpointOptions.uploadUrl}`);
          formData.append(endpointOptions.fileKeyName ?? 'content', file);

          if (endpointOptions.payload) {
            Object.keys(endpointOptions.payload).forEach((key) => {
              formData.append(key, String(endpointOptions.payload?.[key]));
            });
          }
        }

        xhr.send(formData);
        inputRef.current.value = '';
        onProgress?.(refId, 0);
      }
    };

    return (
      <>
        <label htmlFor={name} ref={ref} className={className} {...props} />

        <input
          type='file'
          ref={inputRef}
          accept={accept}
          id={name}
          className='hidden'
          onChange={handleOnChange}
        />
      </>
    );
  },
);

import { useRef, DragEventHandler, ChangeEventHandler } from 'react';

import { useStore } from '@shared/hooks/useStore';

export interface EndpointOptions {
  uploadUrl: string;
  fileKeyName: string;
  downloadUrl?: string;
  payload?: Record<string, string>;
}

export interface FileUploaderProps {
  apiBaseUrl: string;
  onLoadStart?(refId: number): void;
  endpointOptions?: EndpointOptions;
  onChange?(file: File, refId: number): void;
  onLoadEnd?(refId: number, e: unknown): void;
  onError?(refId: number, error: string): void;
  onDragOverChange?(isDragging: boolean): void;
  onLoading?(refId: number, value: boolean): void;
  onProgress?(refId: number, value: number): void;
  onSuccess?(refId: number, response: unknown): void;
}

export const useFileUploader = ({
  onError,
  onChange,
  onLoadEnd,
  onLoading,
  onSuccess,
  onProgress,
  apiBaseUrl,
  onLoadStart,
  endpointOptions,
  onDragOverChange,
}: FileUploaderProps) => {
  const store = useStore();
  const inputRef = useRef<HTMLInputElement>(null);

  const upload = async (file: File) => {
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
        if (xhr.status === 413) {
          onError?.(refId, 'Your file needs to be less than 1MB');

          return;
        }

        onError?.(refId, 'Could not fetch data.');
      }
    };

    if (!endpointOptions) {
      return;
    }

    const formData = new FormData();

    if (endpointOptions) {
      xhr.open(
        'POST',
        `${import.meta.env.VITE_MIDDLEWARE_API_URL}${apiBaseUrl}${
          endpointOptions.uploadUrl
        }`,
      );
      xhr.setRequestHeader(
        'Authorization',
        `Bearer ${store.session.sessionToken}`,
      );
      xhr.setRequestHeader(
        'X-Openline-USERNAME',
        store.session.value.profile.email,
      );
      formData.append(endpointOptions.fileKeyName ?? 'content', file);

      if (endpointOptions.payload) {
        Object.keys(endpointOptions.payload).forEach((key) => {
          formData.append(key, String(endpointOptions.payload?.[key]));
        });
      }
    }

    xhr.send(formData);
    onProgress?.(refId, 0);
  };

  const handleOnChange: ChangeEventHandler<HTMLInputElement> = async () => {
    if (!inputRef) return;

    const currentFiles = inputRef.current?.files;

    if (!currentFiles) return;

    for (const file of currentFiles) {
      await upload(file);
    }

    inputRef.current.value = '';
  };

  const handleDragOver: DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();
    onDragOverChange?.(true);
  };

  const handleDragLeave: DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();
    onDragOverChange?.(false);
  };

  const handleDrop: DragEventHandler<HTMLDivElement> = async (e) => {
    e.preventDefault();
    onDragOverChange?.(false);
    const files = e.dataTransfer?.files;

    if (!files) return;

    for (const file of files) {
      await upload(file);
    }
  };

  return {
    inputRef,
    handleDrop,
    handleDragOver,
    handleOnChange,
    handleDragLeave,
  };
};

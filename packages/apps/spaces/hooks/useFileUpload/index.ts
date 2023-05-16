import { useState } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify';
import { uuid4 } from '@sentry/utils';

export const useFileUpload = ({
  prevFiles = [],
  onBeginFileUpload,
  onFileUpload,
  onFileUploadError,
  onFileRemove,
  uploadInputRef,
}: {
  prevFiles: Array<any>;
  onBeginFileUpload: (data: any) => void;
  onFileUpload: (data: any) => void;
  onFileUploadError: (data: any) => void;
  onFileRemove: (fileId: string) => void;
  uploadInputRef: any;
}) => {
  const [files, setFiles] = useState<any[]>(prevFiles);
  const [isDraggingOver, setIsDraggingOver] = useState(false);

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
    const uploadFileToast = toast.loading('Uploading file');

    axios
      .get(`/fs/jwt`, {
        headers: {
          Accept: 'application/json',
        },
      })
      .then((r: any) => {
        return axios.post(
          `${process.env.FILE_STORAGE_PUBLIC_URL}/file`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              'X-Openline-JWT': r.data.token,
            },
          },
        );
      })
      .then((r: any) => {
        toast.update(uploadFileToast, {
          render: 'File uploaded!',
          type: 'success',
          isLoading: false,
          autoClose: 100,
        });
        clearFileInput();
        onFileUpload({ ...r.data, key: fileKey });

        // return r.data;
      })
      .catch((e) => {
        toast.update(uploadFileToast, {
          render:
            'Oops! We could add this file. Check if file type is supported and can try again or contact our support team',
          type: 'error',
          autoClose: 1000,
          isLoading: false,
        });

        clearFileInput();
        onFileUploadError(fileKey);
      });
  };

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

  const handleFileRemove = (fileKey: any, onFileRemove: any) => {
    setFiles((prevState) => prevState.filter((file) => file.key !== fileKey));
    onFileRemove(fileKey);
  };

  const addFile = (file: any) => {
    setFiles((prevState) => [...prevState, file]);
  };

  return {
    files,
    isDraggingOver,
    handleDrag,
    handleDrop,
    handleInputFileChange,
    handleFileRemove,
    addFile,
  };
};

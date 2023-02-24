import { ChangeEvent } from 'react';
import { Result } from './types';
import axios from 'axios';
import { toast } from 'react-toastify';

interface Props {
  addFileToTextContent: (data: string) => void;
}
export const useFileData = ({ addFileToTextContent }: Props): Result => {
  const handleFetchFile = (id: string) => {
    return fetch(`/fs/file/${id}/download`)
      .then(async (response: any) => {
        const blob = await response.blob();
        const reader = new FileReader();
        reader.onload = function () {
          const dataUrl = reader.result as any;

          if (dataUrl) {
            addFileToTextContent(
              `<img width="400" src='${dataUrl}' alt='${id}'>`,
            );
          } else {
            toast.error('');
          }
        };
        reader.readAsDataURL(blob);
      })
      .catch((reason: any) => {
        toast.error('Oops! We could not load provided file');
      });
  };
  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) {
      return;
    }

    const formData = new FormData();
    formData.append('file', e.target.files[0]);
    axios
      .post(`/fs/file`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      .then((r: any) => handleFetchFile(r.data.id))
      .catch((reason: any) => {
        toast.error(
          'Oops! We could add this file. Check if file type is supported and can try again or contact our support team',
        );
      });
  };

  return {
    onFileChange: handleFileChange,
  };
};

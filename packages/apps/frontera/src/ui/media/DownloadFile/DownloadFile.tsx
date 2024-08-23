import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Button, ButtonProps } from '@ui/form/Button/Button';

interface DownloadFileProps extends ButtonProps {
  fileId: string;
  fileName: string;
  variant: 'outline' | 'solid' | 'ghost' | 'link';
}

export const DownloadFile = observer(
  ({ fileId, variant, fileName, ...rest }: DownloadFileProps) => {
    const { files } = useStore();

    const handleDownload = () => {
      files.downloadAttachment(fileId, fileName);
      files.clear(fileId);
    };

    return (
      <Button variant={variant} onClick={handleDownload} {...rest} size='xs'>
        Download
      </Button>
    );
  },
);

import { forwardRef } from 'react';

import { useFileUploader, FileUploaderProps } from './useFileUploader';

interface FileDropUploaderProps extends FileUploaderProps {
  className?: string;
  children?: React.ReactNode;
  onDragOverChange?: (isDraggingOver: boolean) => void;
}

export const FileDropUploader = forwardRef<
  HTMLDivElement,
  FileDropUploaderProps
>(({ className, children, ...props }, ref) => {
  const { handleDragOver, handleDrop, handleDragLeave } =
    useFileUploader(props);

  return (
    <div
      ref={ref}
      onDrop={handleDrop}
      className={className}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
    >
      {children}
    </div>
  );
});

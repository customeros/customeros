import React from 'react';
import { IconButton } from './IconButton';

interface DeleteIconButtonProps {
  onDelete: () => void;
  style?: any;
}

export const DeleteIconButton: React.FC<DeleteIconButtonProps> = ({
  onDelete,
  style,
}) => {
  return (
    <IconButton
      size={'xxxxs'}
      mode='danger'
      style={{
        width: '12px',
        height: '12px',
        lineHeight: '12px',
        ...style,
      }}
      onClick={onDelete}
      icon={<>-</>}
    />
  );
};

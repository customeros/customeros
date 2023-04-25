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
    <div>
      <IconButton
        size={'xxxxs'}
        mode='danger'
        style={{
          width: '11px',
          height: '11px',
          textAlign: 'center',
          fontSize: '11px',
          // justifyContent: 'center',
          // alignItems: 'center',
          ...style,
        }}
        onClick={onDelete}
        icon={<span>&#8211;</span>}
      />
    </div>
  );
};

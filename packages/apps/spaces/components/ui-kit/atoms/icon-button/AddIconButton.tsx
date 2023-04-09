import React from 'react';
import { IconButton } from './IconButton';
import { Plus } from '../icons';

interface DeleteIconButtonProps {
  onAdd: () => void;
  style?: any;
}

export const AddIconButton: React.FC<DeleteIconButtonProps> = ({
  onAdd,
  style,
}) => {
  return (
    <IconButton
      size={'xxxxs'}
      mode='text'
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        ...style,
      }}
      onClick={onAdd}
      icon={<Plus style={{ transform: 'scale(0.6)' }} />}
    />
  );
};

import React from 'react';
import { IconButton } from './IconButton';

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
      label='Add'
      mode='text'
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        ...style,
      }}
      onClick={onAdd}
      icon={<span>+</span>}
    />
  );
};

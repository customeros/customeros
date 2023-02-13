import React, { FC } from 'react';
import styles from './workspace-button.module.scss';
interface WorkspaceButtonProps {
  label: string;
  image?: string;
  onClick: () => void;
}

export const WorkspaceButton: FC<WorkspaceButtonProps> = ({
  label,
  image,
  onClick,
}) => {
  return (
    <button
      className={styles.button}
      onClick={onClick}
      aria-label={label}
      title={label}
    >
      {image ? (
        <img src={image} alt='' className={styles.imageButton} />
      ) : (
        <span>{label[0]}</span>
      )}
    </button>
  );
};

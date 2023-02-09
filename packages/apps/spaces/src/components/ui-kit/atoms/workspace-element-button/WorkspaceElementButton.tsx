import React from 'react';

interface WorkspaceElementButtonProps {
  onClick: () => void;
  label: string;
  image?: string;
}

export const WorkspaceElementButton: React.FC<WorkspaceElementButtonProps> = ({
  onClick,
  label,
  image,
}) => {
  return (
    <button onClick={onClick}>
      {image && <img src={image} alt={label} />}
      <span>{label}</span>
    </button>
  );
};

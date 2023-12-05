import React from 'react';

import { ContractStatusUpdatedActionPreview } from './contract/ContractStatusUpdatedActionPreview';

interface ActionPreviewModalProps {
  type: string;
}

export const ActionPreviewModal: React.FC<ActionPreviewModalProps> = ({
  type,
}) => {
  switch (type) {
    case 'CONTRACT_STATUS_UPDATED':
      return <ContractStatusUpdatedActionPreview />;
    default:
      return null;
  }
};

import React from 'react';

import { ServiceUpdatedActionPreview } from './service/ServiceUpdatedActionPreview';
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
    case 'SERVICE_LINE_ITEM_QUANTITY_UPDATED':
      return <ServiceUpdatedActionPreview />;
    default:
      return null;
  }
};

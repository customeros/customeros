import React from 'react';

import { ActionType } from '@graphql/types';

import { ServiceUpdatedActionPreview } from './service/ServiceUpdatedActionPreview';
import { ContractStatusUpdatedActionPreview } from './contract/ContractStatusUpdatedActionPreview';

interface ActionPreviewModalProps {
  type: ActionType;
}

export const ActionPreviewModal: React.FC<ActionPreviewModalProps> = ({
  type,
}) => {
  switch (type) {
    case ActionType.ContractStatusUpdated:
      return <ContractStatusUpdatedActionPreview />;
    case ActionType.ServiceLineItemQuantityUpdated:
    case ActionType.ServiceLineItemPriceUpdated:
      return <ServiceUpdatedActionPreview />;
    default:
      return null;
  }
};

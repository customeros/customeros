import React from 'react';
import { RenewalForecastUpdatedActionPreview } from './renewal-forecast/RenewalForecastUpdatedActionPreview';
import { RenewalLikelihoodUpdatedActionPreview } from './renewal-likelihood/RenewalLikelihoodUpdatedActionPreview';

interface ActionPreviewModalProps {
  type: string;
}

export const ActionPreviewModal: React.FC<ActionPreviewModalProps> = ({
  type,
}) => {
  switch (type) {
    case 'RENEWAL_FORECAST_UPDATED':
      return <RenewalForecastUpdatedActionPreview />;
    case 'RENEWAL_LIKELIHOOD_UPDATED':
      return <RenewalLikelihoodUpdatedActionPreview />;
    default:
      return null;
  }
};

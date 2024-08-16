import { ActionType } from '@graphql/types';
import { InvoicePreviewModal } from '@organization/components/Timeline/PastZone/events/invoice/InvoicePreviewModal';

import { ServiceUpdatedActionPreview } from './service/ServiceUpdatedActionPreview';
import { ContractStatusUpdatedActionPreview } from './contract/ContractStatusUpdatedActionPreview';
import { OnboardingStatusChangedActionPreview } from './onboarding/OnboardingStatusChangedActionPreview';

interface ActionPreviewModalProps {
  type: ActionType;
}

export const ActionPreviewModal = ({ type }: ActionPreviewModalProps) => {
  switch (type) {
    case ActionType.ContractStatusUpdated:
    case ActionType.ContractRenewed:
      return <ContractStatusUpdatedActionPreview />;
    case ActionType.ServiceLineItemQuantityUpdated:
    case ActionType.ServiceLineItemPriceUpdated:
    case ActionType.ServiceLineItemBilledTypeUpdated:
      return <ServiceUpdatedActionPreview />;
    case ActionType.ServiceLineItemBilledTypeOnceCreated:
    case ActionType.ServiceLineItemBilledTypeUsageCreated:
    case ActionType.ServiceLineItemBilledTypeRecurringCreated:
      return <ServiceUpdatedActionPreview mode='created' />;
    case ActionType.OnboardingStatusChanged:
      return <OnboardingStatusChangedActionPreview />;
    default:
      return <InvoicePreviewModal />;
  }
};

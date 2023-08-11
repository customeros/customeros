'use client';
import { useState } from 'react';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { BillingDetailsCard } from '@organization/components/Tabs/panels/AccountPanel/components/BillingDetailsCard/BillingDetailsCard';
import {
  RenewalLikelihood,
  Value as RenewalLikelihoodValue,
} from './RenewalLikelihood';
import {
  RenewalForecast,
  Value as RenewalForecastValue,
} from './RenewalForecast';

export const AccountPanel = () => {
  const [renewalLikelihood, setRenewalLikelihood] =
    useState<RenewalLikelihoodValue>({ reason: '', likelihood: 'NOT_SET' });
  const [renewalForecast, setRenewalForecast] = useState<RenewalForecastValue>({
    reason: '',
    forecast: '',
  });

  return (
    <OrganizationPanel title='Account'>
      <RenewalLikelihood
        value={renewalLikelihood}
        onChange={setRenewalLikelihood}
      />
      <RenewalForecast value={renewalForecast} onChange={setRenewalForecast} />
      <BillingDetailsCard />
    </OrganizationPanel>
  );
};

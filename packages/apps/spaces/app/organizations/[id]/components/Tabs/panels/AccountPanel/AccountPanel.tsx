'use client';
import { useState } from 'react';

import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

import {
  RenewalLikelihood,
  Value as RenewalLikelihoodValue,
} from './RenewalLikelihood';

export const AccountPanel = () => {
  const [renewalLikelihood, setRenewalLikelihood] =
    useState<RenewalLikelihoodValue>({ reason: '', likelihood: 'NOT_SET' });

  return (
    <OrganizationPanel title='Account'>
      <RenewalLikelihood
        value={renewalLikelihood}
        onChange={setRenewalLikelihood}
      />
    </OrganizationPanel>
  );
};

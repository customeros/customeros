import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ContractCard } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractCard';
import { ARRForecast } from '@organization/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast.tsx';
import { ContractModalsContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext.tsx';
import { ContractModalStatusContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext.tsx';

export const Contracts = observer(() => {
  const id = useParams()?.id as string;
  const store = useStore();
  const organization = store.organizations.getById(id);
  const contracts = organization?.contracts;

  return (
    <>
      <ARRForecast
        name={organization?.value.name || ''}
        currency={contracts?.[0]?.currency || 'USD'}
        arrForecast={organization?.value?.renewalSummaryArrForecast}
        maxArrForecast={organization?.value?.renewalSummaryMaxArrForecast}
        renewalLikelihood={organization?.value?.renewalSummaryRenewalLikelihood}
      />
      {contracts?.map((c) => {
        return (
          <div
            key={`contract-card-${c.metadata.id}`}
            className='flex gap-4 flex-col w-full mb-4'
          >
            <ContractModalStatusContextProvider id={c.metadata.id}>
              <ContractModalsContextProvider id={c.metadata.id}>
                <ContractCard
                  values={c}
                  organizationName={organization?.value.name || ''}
                />
              </ContractModalsContextProvider>
            </ContractModalStatusContextProvider>
          </div>
        );
      })}
    </>
  );
});

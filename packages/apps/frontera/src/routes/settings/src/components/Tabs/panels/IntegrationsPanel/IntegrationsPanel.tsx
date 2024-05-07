import { useState, useEffect } from 'react';

import Fuse from 'fuse.js';
import { autorun } from 'mobx';
import { observer } from 'mobx-react-lite';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';

import { Input } from '@ui/form/Input/Input';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';

import { IntegrationItem, integrationsData } from './data';
import { SettingsIntegrationItem } from './SettingsIntegrationItem';

export const IntegrationsPanel = observer(() => {
  const { settingsStore } = useStore();
  const iApp = useIntegrationApp();
  const { items: iIntegrations } = useIntegrations();
  const { items: iConnections, refresh } = useConnections();

  const [integrations, setIntegrations] =
    useState<IntegrationItem[]>(integrationsData);

  const [integrationsDisplayed, setIntegrationsDisplayed] = useState<
    IntegrationItem[]
  >([]);

  useEffect(() => {
    autorun(() => {
      const map = integrations.map((integration) => {
        return {
          ...integration,
          state:
            settingsStore.integrations.value[integration.key]?.state ??
            'INACTIVE',
        };
      });
      setIntegrations(map);
      setIntegrationsDisplayed(map);
    });
  }, []);

  const handleFilterResults = (value: string) => {
    if (value.length === 0) {
      setIntegrationsDisplayed(integrations);

      return;
    }

    // Options for Fuse
    const options = {
      keys: ['key'],
      shouldSort: true,
      caseSensitive: false,
      includeScore: true,
      findAllMatches: true,
    };

    const fuse = new Fuse(integrations, options);
    const result = fuse.search(value);
    const finalResult = result.map((res) => res.item);
    setIntegrationsDisplayed(finalResult);
  };

  // integration.app related logic (temporary)
  const activeIntegrations = iConnections.map((item) => item.integration?.key);
  const availableIntegrations = iIntegrations.map((item) => item.key);

  const handleIntegration = (integrationKey: string) => async () => {
    const option = iIntegrations.find((item) => item.key === integrationKey);

    if (!option) {
      return;
    }
    try {
      await iApp.integration(option.key).open({ showPoweredBy: false });
      await refresh();
    } catch (err) {
      toastError('Integration failed', 'get-intergration-data');
    }
  };

  return (
    <>
      <div className=' flex h-[calc(100vh-1rem)] max-w-[600px] bg-gray-25  rounded-2xl flex-col max-h-[calc(100vh - 1rem)] relative '>
        <div className='pb-1 pt-5 px-6 '>
          <h1 className='text-2xl font-bold'>Data Integrations</h1>
          <Input
            onChange={(event) => handleFilterResults(event.target.value)}
            placeholder={'Search...'}
          />
        </div>
        <div className='overflow-auto pt-1 px-5 pb-5 w-full'>
          <h3 className='text-lg font-medium'>Active integrations</h3>
          {settingsStore.integrations.isLoading && (
            <div className='flex-col space-y-3 my-2'>
              <Skeleton className='h-5 w-full rounded-sm' />
              <Skeleton className='h-5 w-full rounded-sm' />
            </div>
          )}
          {!settingsStore.integrations.isLoading && (
            <>
              {integrationsDisplayed
                .filter((integration: IntegrationItem) => {
                  if (integration.isFromIntegrationApp) {
                    return activeIntegrations.includes(integration.key);
                  } else {
                    return integration.state === 'ACTIVE';
                  }
                })
                .map((integration: IntegrationItem) => {
                  const option = integration.key;
                  const isFromIApp = activeIntegrations.includes(option);

                  return (
                    <SettingsIntegrationItem
                      key={integration.key}
                      icon={integration.icon}
                      identifier={integration.identifier}
                      name={integration.name}
                      onDisable={
                        isFromIApp ? handleIntegration(option) : undefined
                      }
                      state={isFromIApp ? 'ACTIVE' : integration.state}
                      fields={integration.fields}
                    />
                  );
                })}

              {!integrationsDisplayed.filter(
                (integration: IntegrationItem) =>
                  integration.state === 'ACTIVE',
              ).length && (
                <p className='text-gray-400 mt-1 mb-3'>
                  There are no active integrations
                </p>
              )}
            </>
          )}

          <h3 className='text-lg font-medium'>Inactive integrations</h3>
          {settingsStore.integrations.isLoading && (
            <div className='flex-col space-y-3 mt-2'>
              <Skeleton className='h-5 w-full rounded-sm' />
              <Skeleton className='h-5 w-full rounded-sm' />
              <Skeleton className='h-5 w-full rounded-sm' />
            </div>
          )}
          {!settingsStore.integrations.isLoading && (
            <>
              {integrationsDisplayed
                .filter((integration: IntegrationItem) => {
                  if (integration.isFromIntegrationApp) {
                    return !activeIntegrations.includes(integration.key);
                  } else {
                    return integration.state === 'INACTIVE';
                  }
                })
                .map((integration: IntegrationItem) => {
                  const option = integration.key;
                  const isFromIApp = availableIntegrations.includes(option);

                  return (
                    <SettingsIntegrationItem
                      key={integration.key}
                      icon={integration.icon}
                      identifier={integration.identifier}
                      name={integration.name}
                      state={integration.state}
                      onEnable={
                        isFromIApp ? handleIntegration(option) : undefined
                      }
                      fields={integration.fields}
                    />
                  );
                })}
            </>
          )}
        </div>
      </div>
    </>
  );
});

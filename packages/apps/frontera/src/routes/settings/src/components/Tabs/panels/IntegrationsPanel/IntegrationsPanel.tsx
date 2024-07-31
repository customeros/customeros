import { useMemo, useEffect } from 'react';

import Fuse from 'fuse.js';
import { autorun } from 'mobx';
import { observer, useLocalObservable } from 'mobx-react-lite';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';

import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';

import { IntegrationItem, integrationsData } from './data';
import { SettingsIntegrationItem } from './SettingsIntegrationItem';

export const IntegrationsPanel = observer(() => {
  const iApp = useIntegrationApp();
  const store = useStore();
  const { items: iIntegrations } = useIntegrations();
  const { items: iConnections, refresh } = useConnections();

  // integration.app related logic (temporary)
  const iAppConnections = iConnections
    .map((item) => item.integration?.key)
    .filter(Boolean);

  const state = useLocalObservable(() => ({
    all: integrationsData.map((integration) => ({
      ...integration,
      state:
        store.settings.integrations.value[integration.key]?.state ?? 'INACTIVE',
    })),
    setAll(integrations: IntegrationItem[]) {
      this.all = integrations;
    },
    searchTerm: '',
    setSearchTerm(term: string) {
      this.searchTerm = term;
    },
  }));

  const mixedActiveIntegrations = useMemo(
    () =>
      state.all.filter((integration) => {
        if (integration.isFromIntegrationApp) {
          return iAppConnections.includes(integration.key);
        }

        return integration.state === 'ACTIVE';
      }),
    [state.all, iAppConnections],
  );

  const mixedInactiveIntegrations = useMemo(
    () =>
      state.all.filter((integration) => {
        if (integration.isFromIntegrationApp) {
          return !iAppConnections.includes(integration.key);
        }

        return integration.state === 'INACTIVE';
      }),
    [state.all, iAppConnections],
  );

  const searchedIntegrations = useMemo(() => {
    if (state.searchTerm === '') return mixedInactiveIntegrations;

    const options = {
      keys: ['key'],
      shouldSort: true,
      caseSensitive: false,
      includeScore: true,
      findAllMatches: true,
    };

    const fuse = new Fuse(mixedInactiveIntegrations, options);
    const result = fuse.search(state.searchTerm);

    return result.map((res) => res.item);
  }, [state.searchTerm, mixedInactiveIntegrations]);

  const handleIntegration = (integrationKey: string) => async () => {
    const option = iIntegrations.find((item) => item.key === integrationKey);

    if (!option) {
      return;
    }

    try {
      await iApp.integration(option.key).open({ showPoweredBy: false });
      await refresh();
    } catch (err) {
      toastError('Integration failed', 'get-integration-data');
    }
  };

  useEffect(() => {
    const dispose = autorun(() => {
      state.setAll(
        integrationsData.map((integration) => ({
          ...integration,
          state:
            store.settings.integrations.value[integration.key]?.state ??
            'INACTIVE',
        })),
      );
    });

    return () => {
      dispose();
    };
  }, []);

  return (
    <>
      <div className='flex h-[calc(100vh-1rem)] max-w-[600px] bg-gray-25 rounded-2xl flex-col max-h-[calc(100vh - 1rem)] relative'>
        <div className='pb-1 pt-5 px-6'>
          <h1 className='text-2xl font-bold'>Data Integrations</h1>
          <Input
            value={state.searchTerm}
            placeholder={'Search...'}
            onChange={(event) => state.setSearchTerm(event.target.value)}
          />
        </div>
        <div className='overflow-auto pt-1 px-5 pb-5 w-full'>
          <h3 className='text-lg font-medium'>Active integrations</h3>
          {mixedActiveIntegrations.map((integration) => {
            const option = integration.key;
            const isFromIApp = iAppConnections.includes(option);

            return (
              <SettingsIntegrationItem
                key={integration.key}
                icon={integration.icon}
                name={integration.name}
                fields={integration.fields}
                isIntegrationApp={isFromIApp}
                identifier={integration.identifier}
                onSuccess={() => state.setSearchTerm('')}
                state={isFromIApp ? 'ACTIVE' : integration.state}
                onDisable={isFromIApp ? handleIntegration(option) : undefined}
              />
            );
          })}

          {!mixedActiveIntegrations.length && (
            <p className='text-gray-400 mt-1 mb-3'>
              There are no active integrations
            </p>
          )}

          <h3 className='text-lg font-medium mt-4'>Inactive integrations</h3>
          {searchedIntegrations.map((integration) => {
            const option = integration.key;
            const isFromIApp = integration.isFromIntegrationApp;

            return (
              <SettingsIntegrationItem
                key={integration.key}
                icon={integration.icon}
                name={integration.name}
                state={integration.state}
                fields={integration.fields}
                isIntegrationApp={isFromIApp}
                identifier={integration.identifier}
                onSuccess={() => state.setSearchTerm('')}
                onEnable={isFromIApp ? handleIntegration(option) : undefined}
              />
            );
          })}
        </div>
      </div>
    </>
  );
});

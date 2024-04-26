'use client';
import React, { useRef, useState, useEffect } from 'react';

import Fuse from 'fuse.js';
import { GetIntegrationsSettings } from 'services';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';

import { Input } from '@ui/form/Input/Input';
import { Skeleton } from '@ui/feedback/Skeleton';
import { toastError } from '@ui/presentation/Toast';

import { IntegrationItem, integrationsData } from './data';
import { SettingsIntegrationItem } from './SettingsIntegrationItem';

export const IntegrationsPanel = () => {
  const iApp = useIntegrationApp();
  const { items: iIntegrations } = useIntegrations();
  const { items: iConnections, refresh } = useConnections();
  const [reload, setReload] = useState<boolean>(true);
  const reloadRef = useRef<boolean>(reload);

  const [loading, setLoading] = useState<boolean>(true);

  const [integrations, setIntegrations] =
    useState<IntegrationItem[]>(integrationsData);

  const [integrationsDisplayed, setIntegrationsDisplayed] = useState<
    IntegrationItem[]
  >([]);

  useEffect(() => {
    GetIntegrationsSettings()
      .then((data) => {
        const map = integrations.map((integration) => {
          return {
            ...integration,
            state:
              (data as Record<string, IntegrationItem>)[integration.key]
                ?.state ?? 'INACTIVE',
          };
        });

        setIntegrations(map);
        setIntegrationsDisplayed(map);

        setLoading(false);
      })
      .catch(() => {
        setLoading(false);

        toastError(
          'There was a problem on our side and we cannot load settings data at the moment,  we are doing our best to solve it! ',
          'get-intergration-data',
        );
      });
  }, [reload]);

  const handleFilterResults = (value: string) => {
    if (value.length === 0) {
      setIntegrationsDisplayed(integrations);

      return;
    }

    // Options for Fuse
    const options = {
      // which keys to search in
      keys: ['key'],
      // turn on case sensitivity
      shouldSort: true,
      // specify whether comparisons should be case sensitive
      caseSensitive: false,
      includeScore: true, // doesn't have to be true, it's just an example
      findAllMatches: true, // doesn't have to be true, it's just an example
    };

    const fuse = new Fuse(integrations, options);

    const result = fuse.search(value);

    // If you want only the original list items and in array format, you can map over the results:
    const finalResult = result.map((res) => res.item);

    // Update the display
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
      <div className=' flex h-[calc(100vh-1rem)] bg-gray-25  rounded-2xl flex-col max-w-[50%] max-h-[calc(100vh - 1rem)] relative '>
        <div className='pb-1 pt-5 px-6 '>
          <h1 className='text-2xl font-bold'>Data Integrations</h1>
          <Input
            onChange={(event) => handleFilterResults(event.target.value)}
            placeholder={'Search...'}
          />
        </div>
        <div className='overflow-auto pt-1 px-5 pb-5'>
          <h3 className='text-lg font-medium'>Active integrations</h3>
          {loading && (
            <div className='flex-col space-y-3 my-2'>
              <Skeleton className='h-5 w-full rounded-sm' />
              <Skeleton className='h-5 w-full rounded-sm' />
            </div>
          )}
          {!loading && (
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
                      settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                      }}
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
          {loading && (
            <div className='flex-col space-y-3 mt-2'>
              <Skeleton className='h-5 w-full rounded-sm' />
              <Skeleton className='h-5 w-full rounded-sm' />
              <Skeleton className='h-5 w-full rounded-sm' />
            </div>
          )}
          {!loading && (
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
                      settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                      }}
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
};

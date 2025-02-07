import { createPortal } from 'react-dom';
import { useMemo, useEffect } from 'react';

import { toJS } from 'mobx';
import get from 'lodash/get';
import { useKey } from 'rooks';
import { match } from 'ts-pattern';
import { Tracer } from '@infra/tracer';
import { type RootStore } from '@store/root';
import { Observer, observer } from 'mobx-react-lite';
import { CommonStore } from '@store/Common/Common.store';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { Switch } from '@ui/form/Switch';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Resizable } from '@ui/presentation/Resizable';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

import { Traces } from './Traces';
import { DevtoolsStore } from './state';
import { useStore } from '../../hooks/useStore';

const devTools = new DevtoolsStore();

type StoreReturnType =
  | RootStore['organizations']
  | RootStore['tableViewDefs']
  | RootStore['contacts']
  | RootStore['contracts']
  | RootStore['flows']
  | RootStore['jobRoles']
  | RootStore['tags']
  | RootStore['agents']
  | RootStore['common']
  | null;

interface DevtoolsProps {
  open: boolean;
  onClose: () => void;
  onToggle: () => void;
}

export const Devtools = observer(
  ({ open, onClose, onToggle }: DevtoolsProps) => {
    const store = useStore();
    const hasFlagEnabled = useFeatureIsOn('devtools');

    const ENABLED = useMemo(
      () => import.meta.env.DEV || hasFlagEnabled,
      [hasFlagEnabled],
    );

    const defaultWidht = window.innerWidth / 2;
    const defaultX = window.innerWidth - defaultWidht * 1.5;
    const defaultY = window.innerHeight / 3;

    useKey('`', onToggle, { when: ENABLED });

    useEffect(() => {
      if (!ENABLED) return;

      const handleGqlReq = (e: unknown) => {
        const reqId = get(e, 'detail.reqId');
        const reqName = get(e, 'detail.name');
        const reqVariables = get(e, 'detail.variables');

        devTools.addGqlOp({
          id: reqId ?? crypto.randomUUID(),
          name: reqName ?? crypto.randomUUID(),
          variables: reqVariables ?? null,
        });
      };

      const handleGqlRes = (e: unknown) => {
        const reqId = get(e, 'detail.reqId');
        const resData = get(e, 'detail.data');
        const resErrors = get(e, 'detail.errors');

        if (!reqId) return;

        devTools.addGqlRes({
          id: reqId,
          data: resData ?? null,
          errors: resErrors ?? null,
        });
      };

      window.addEventListener('gql-req', handleGqlReq);
      window.addEventListener('gql-res', handleGqlRes);

      return () => {
        if (!ENABLED) return;
        window.removeEventListener('gql-req', handleGqlReq);
        window.removeEventListener('gql-res', handleGqlRes);
      };
    }, []);

    const detailedGqlOperation = devTools.gqlOperations.find(
      (o) => o.id === devTools.openGqlOperationId,
    );
    const detailedGqlResponse = devTools.gqlResponses.get(
      devTools.openGqlOperationId ?? '',
    );
    const detailedStore = match(devTools.detailedStore)
      .returnType<StoreReturnType>()
      .with('tableViewDefs', () => store.tableViewDefs)
      .with('organizations', () => store.organizations)
      .with('contacts', () => store.contacts)
      .with('contracts', () => store.contracts)
      .with('flows', () => store.flows)
      .with('jobRoles', () => store.jobRoles)
      .with('tags', () => store.tags)
      .with('agents', () => store.agents)
      .with('common.slackChannels', () => store.common)
      .otherwise(() => null);

    const getStoreValue = (store: StoreReturnType) => {
      if (store instanceof CommonStore) {
        return store.slackChannels;
      }

      return store?.value;
    };

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const getEntityValue = (entity: Record<string, any>) =>
      match(devTools.detailedStore)
        .with('organizations', () => get(entity, 'value', {}))
        .with('tableViewDefs', () => get(entity, 'value', {}))
        .with('contacts', () => get(entity, 'value', {}))
        .with('contracts', () => get(entity, 'value', {}))
        .with('flows', () => get(entity, 'value', {}))
        .with('jobRoles', () => get(entity, 'value', {}))
        .with('tags', () => get(entity, 'value', {}))
        .with('agents', () => get(entity, 'value', {}))
        .with('common.slackChannels', () => entity)
        .otherwise(() => entity);

    const getEntityName = (entity: Record<string, string>) =>
      match(devTools.detailedStore)
        .returnType<string>()
        .with('organizations', () => get(entity, 'value.name', 'Unnamed'))
        .with('tableViewDefs', () => get(entity, 'name', 'Unnamed'))
        .with('contacts', () => get(entity, 'name', 'Unnamed'))
        .with('contracts', () => get(entity, 'value.contractName', 'Unnamed'))
        .with('flows', () => get(entity, 'value.name', 'Unnamed'))
        .with('jobRoles', () => get(entity, 'value.jobTitle', 'Unnamed'))
        .with('tags', () => get(entity, 'value.name', 'Unnamed'))
        .with('agents', () => get(entity, 'value.name', 'Unnamed'))
        .with('common.slackChannels', () => get(entity, 'name', 'Unnamed'))
        .otherwise(() => 'Unnamed');

    if (!open) return null;

    return createPortal(
      <Resizable
        defaultHeight={350}
        defaultWidth={window.innerWidth / 2}
        className='absolute bottom-0 left-0 right-0'
        defaultPosition={{
          x: defaultX,
          y: defaultY,
        }}
      >
        {(isDragging, startMove) => (
          <Observer>
            {() => (
              <div
                className={cn(
                  'relative flex bg-white drop-shadow-2xl ring-1 ring-gray-200 w-full h-full z-10 flex-col p-0 rounded-lg transition',
                  isDragging && 'ring-2 ring-primary-400',
                )}
              >
                <IconButton
                  size='xxs'
                  icon={<X />}
                  variant='ghost'
                  onClick={onClose}
                  aria-label='close devtools'
                  className='absolute top-1 right-1'
                />

                <div
                  onMouseDown={(e) => startMove(e, true)}
                  className='border-b rounded-t-lg border-b-gray-200 pb-0.5 pl-0.5 w-full bg-gray-50 hover:cursor-grab'
                >
                  <ButtonGroup variant='old'>
                    <Button
                      size='xxs'
                      onClick={() => devTools.toggleView('operations')}
                      className={cn(
                        devTools.view === 'operations' &&
                          'bg-primary-100 focus:bg-primary-200 hover:bg-primary-200',
                      )}
                    >
                      Operations
                    </Button>
                    <Button
                      size='xxs'
                      onClick={() => devTools.toggleView('store')}
                      className={cn(
                        devTools.view === 'store' &&
                          'bg-primary-100 focus:bg-primary-200 hover:bg-primary-200',
                      )}
                    >
                      Store
                    </Button>
                    <Button
                      size='xxs'
                      onClick={() => devTools.toggleView('traces')}
                      className={cn(
                        devTools.view === 'traces' &&
                          'bg-primary-100 focus:bg-primary-200 hover:bg-primary-200',
                      )}
                    >
                      Traces
                    </Button>
                    <Button
                      size='xxs'
                      onClick={() => devTools.toggleView('settings')}
                      className={cn(
                        devTools.view === 'settings' &&
                          'bg-primary-100 focus:bg-primary-200 hover:bg-primary-200',
                      )}
                    >
                      Settings
                    </Button>
                  </ButtonGroup>
                </div>
                <div className='flex h-full overflow-hidden'>
                  <Resizable
                    isMovable={false}
                    defaultWidth={200}
                    resizeDirection='horizontal'
                    className='border-r border-r-gray-200'
                  >
                    <ScrollAreaRoot className='h-full w-full overflow-hidden'>
                      <ScrollAreaViewport className='h-full'>
                        <div className='bg-white border-b border-b-gray-200 p-1'>
                          <Input
                            size='xs'
                            variant='outline'
                            className='text-xs'
                            placeholder='Search...'
                            value={devTools.operationsSearchTerm}
                            onChange={(e) =>
                              devTools.searchOperations(e.target.value)
                            }
                          />
                        </div>

                        <div className={cn('flex flex-col h-full')}>
                          {devTools.view === 'operations' &&
                            devTools.filteredGqlOperations.map((op) => {
                              const isSelected =
                                detailedGqlOperation?.id === op.id;
                              const hasError = devTools.gqlResponses.get(
                                op.id ?? '',
                              )?.errors;

                              return (
                                <div
                                  key={op.id}
                                  onClick={() => {
                                    devTools.toggleGqlOp(op.id);
                                  }}
                                  className={cn(
                                    'border-b border-b-gray-200 cursor-pointer hover:bg-gray-50 pl-1',
                                    isSelected && !hasError && 'bg-gray-100',
                                    isSelected && hasError && 'bg-error-100',
                                    !isSelected &&
                                      hasError &&
                                      'bg-error-50 hover:bg-error-100',
                                  )}
                                >
                                  <span
                                    className={cn(
                                      'text-xs leading-0',
                                      hasError && 'text-error-500',
                                      isSelected && 'font-medium',
                                    )}
                                  >
                                    {op.name}
                                  </span>
                                </div>
                              );
                            })}
                          {devTools.view === 'store' &&
                            devTools.visibleStores.map((name) => {
                              const isSelected =
                                devTools.detailedStore === name;

                              return (
                                <div
                                  key={name}
                                  onClick={() => {
                                    devTools.toggleStore(name);
                                  }}
                                  className={cn(
                                    'border-b border-b-gray-200 cursor-pointer hover:bg-gray-50 pl-1',
                                    isSelected && 'bg-gray-100',
                                  )}
                                >
                                  <span
                                    className={cn(
                                      'text-xs leading-0',
                                      isSelected && 'font-medium',
                                    )}
                                  >
                                    {name}
                                  </span>
                                </div>
                              );
                            })}
                        </div>
                      </ScrollAreaViewport>
                      <ScrollAreaScrollbar orientation='vertical'>
                        <ScrollAreaThumb />
                      </ScrollAreaScrollbar>
                    </ScrollAreaRoot>
                  </Resizable>

                  <ScrollAreaRoot className='h-full w-full overflow-hidden p-1'>
                    <ScrollAreaViewport>
                      {detailedGqlResponse &&
                        devTools.view === 'operations' && (
                          <>
                            <Input
                              size='xs'
                              variant='outline'
                              className='text-xs'
                              placeholder='Search...'
                              value={devTools.operationsSearchTerm}
                              onChange={(e) =>
                                devTools.searchOperations(e.target.value)
                              }
                            />
                            <div className='flex flex-col space-y-1 w-[100px] pb-4'>
                              <p className='text-sm font-medium'>
                                {detailedGqlOperation?.name}
                              </p>
                              <pre className='text-xs'>
                                variables:{' '}
                                {JSON.stringify(
                                  detailedGqlOperation?.variables,
                                  null,
                                  2,
                                )}
                              </pre>
                              <pre className='text-xs'>
                                response:{' '}
                                {JSON.stringify(
                                  detailedGqlResponse?.data,
                                  null,
                                  2,
                                )}
                              </pre>

                              <pre className='text-xs text-error-500'>
                                errors:{' '}
                                {JSON.stringify(
                                  detailedGqlResponse?.errors,
                                  null,
                                  2,
                                )}
                              </pre>
                            </div>
                          </>
                        )}
                      {devTools.detailedStore && devTools.view === 'store' && (
                        <>
                          <p className='font-medium underline capitalize mb-1'>
                            {devTools.detailedStore}
                          </p>

                          {detailedStore &&
                            // @ts-expect-error - TS is working against us here
                            Array.from(getStoreValue(detailedStore))?.map(
                              ([k, v]) => {
                                const isSelected =
                                  devTools.detailedEntityId === k;
                                const name =
                                  getEntityName(
                                    v as unknown as Record<string, string>,
                                  ) || 'Unnamed';

                                return (
                                  <div
                                    key={k}
                                    className='flex flex-col space-y-1'
                                  >
                                    <div
                                      onClick={() => devTools.toggleEntity(k)}
                                      className={cn(
                                        'flex items-center border-b borde-b-gray-200 cursor-pointer hover:bg-gray-50 py-0.5',
                                        isSelected && 'bg-gray-100',
                                      )}
                                    >
                                      <span className='text-xs font-medium mr-0.5'>
                                        {name}
                                      </span>
                                      <span className='text-xs text-gray-500'>
                                        ({k})
                                      </span>
                                    </div>
                                    {isSelected && (
                                      <pre className='text-xs'>
                                        {JSON.stringify(
                                          toJS(getEntityValue(v)),
                                          null,
                                          2,
                                        )}
                                      </pre>
                                    )}
                                  </div>
                                );
                              },
                            )}
                        </>
                      )}
                      {devTools.view === 'traces' && <Traces />}
                      {devTools.view === 'settings' && (
                        <div className='w-full'>
                          <p className='font-medium underline mb-2'>Settings</p>
                          <div className='flex items-center justify-between w-full'>
                            <span className='text-sm font-medium'>Tracer</span>

                            <div className='flex items-center gap-2 mr-2'>
                              <span
                                className={cn(
                                  'text-xs text-grayModern-500',
                                  Tracer.enabled && 'text-primary-500',
                                )}
                              >
                                {Tracer.enabled ? 'enabled' : 'disabled'}
                              </span>
                              <Switch
                                checked={Tracer.enabled}
                                onChange={() =>
                                  Tracer.enabled
                                    ? Tracer.disable()
                                    : Tracer.enable()
                                }
                              />
                            </div>
                          </div>

                          <div className='flex items-center justify-between w-full pl-4 mt-2'>
                            <span className='text-xs font-medium'>
                              Show Callstack
                            </span>

                            <div className='flex items-center gap-2 mr-2'>
                              <span
                                className={cn(
                                  'text-xs text-grayModern-500',
                                  Tracer.displayCallStack && 'text-primary-500',
                                )}
                              >
                                {Tracer.displayCallStack
                                  ? 'enabled'
                                  : 'disabled'}
                              </span>
                              <Switch
                                checked={Tracer.displayCallStack}
                                onChange={() =>
                                  Tracer.displayCallStack
                                    ? Tracer.hideCallstack()
                                    : Tracer.showCallstack()
                                }
                              />
                            </div>
                          </div>
                        </div>
                      )}
                    </ScrollAreaViewport>
                    <ScrollAreaScrollbar orientation='vertical'>
                      <ScrollAreaThumb />
                    </ScrollAreaScrollbar>
                    <ScrollAreaScrollbar orientation='horizontal'>
                      <ScrollAreaThumb />
                    </ScrollAreaScrollbar>
                  </ScrollAreaRoot>
                </div>
              </div>
            )}
          </Observer>
        )}
      </Resizable>,
      document.body,
    );
  },
);

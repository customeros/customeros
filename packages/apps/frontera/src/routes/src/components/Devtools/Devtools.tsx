import { useEffect } from 'react';
import { createPortal } from 'react-dom';

import { toJS } from 'mobx';
import get from 'lodash/get';
import { useKey } from 'rooks';
import { match } from 'ts-pattern';
import { Observer, observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Resizable } from '@ui/presentation/Resizable';
import { Settings01 } from '@ui/media/icons/Settings01';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

import { DevtoolsStore } from './state';
import { useStore } from '../../hooks/useStore';

const devTools = new DevtoolsStore();

export const Devtools = observer(() => {
  const store = useStore();
  const { open, onOpen, onClose, onToggle } = useDisclosure();

  const defaultWidht = window.innerWidth / 2;
  const defaultX = window.innerWidth - defaultWidht * 1.5;
  const defaultY = window.innerHeight / 3;

  useKey('`', onToggle);

  useEffect(() => {
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
    .returnType<
      (typeof store)['organizations'] | (typeof store)['tableViewDefs'] | null
    >()
    .with('tableViewDefs', () => store.tableViewDefs)
    .with('organizations', () => store.organizations)
    .otherwise(() => null);

  return createPortal(
    open ? (
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
                  <ButtonGroup>
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
                            Array.from(detailedStore?.value)?.map(([k, v]) => {
                              const isSelected =
                                devTools.detailedEntityId === k;

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
                                      {v?.name ?? 'Unnamed'}
                                    </span>
                                    <span className='text-xs text-gray-500'>
                                      ({k})
                                    </span>
                                  </div>
                                  {isSelected && (
                                    <pre className='text-xs'>
                                      {JSON.stringify(toJS(v?.value), null, 2)}
                                    </pre>
                                  )}
                                </div>
                              );
                            })}
                        </>
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
      </Resizable>
    ) : (
      <div className='absolute right-2 bottom-2'>
        <IconButton
          size='lg'
          onClick={onOpen}
          icon={<Settings01 />}
          className='rounded-2xl'
          aria-label='frontera devtools'
        />
      </div>
    ),
    document.body,
  );
});

import React, { useEffect, useRef, useState } from 'react';
import { OverlayPanel } from 'primereact/overlaypanel';
import { Skeleton } from 'primereact/skeleton';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import styles from './search.module.scss';
import classNames from 'classnames';
import { Button } from '../button';
import { Input } from '../input';
interface SearchWithOverlayProps {
  resourceLabel: string;
  overlayWidth: string;
  triggerType: 'dropdown' | 'button';

  //search with icon-button
  buttonLabel: string;
  buttonIcon?: any;

  searchDelay: number;
  searchBy: Array<{
    label: string;
    field: string;
    operation?: string;
  }>;

  //search with dropdown
  value: string;

  searchData: any;
  itemTemplate: any;
  maxResults: number;
  onItemSelected: any;
  options: Array<any>;
  loadingOptions: boolean;
}

export const SearchWithOverlay = ({
  overlayWidth = '400px',
  triggerType = 'dropdown',
  searchDelay = 1000,
  searchBy = [],
  maxResults = 25,
  value = 'none',
  options = [],
  loadingOptions = false,
  ...props
}: Partial<SearchWithOverlayProps>) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const overlayRef = useRef<OverlayPanel>(null);
  const labelRef = useRef<HTMLSpanElement>(null);

  // classess used to make the element active
  const initialContainerClassName = 'p-dropdown w-full flex';
  const selectedContainerClassName = ' p-focus';
  const [currentContainerClassName, setCurrentContainerClassName] = useState(
    initialContainerClassName,
  );

  const initialSelectedItemClassName = ' p-dropdown-label p-inputtext w-full';
  const selectedItemEmptyClassName = ' p-dropdown-label-empty';
  const [selectedItemClassName, setSelectedItemClassName] = useState(
    initialSelectedItemClassName +
      (!value || value === '' ? selectedItemEmptyClassName : ''),
  );

  const [loadingData, setLoadingData] = useState(false);

  const [searchResultList, setSearchResultList] = useState([] as any);
  const [totalElements, setTotalElements] = useState([] as any);

  const [displayValue, setDisplayValue] = useState(value);
  const [filters, setFilters] = useState(
    searchBy.map((sb: any) => {
      return {
        label: sb.label,
        field: sb.field,
        value: undefined,
        operation: sb.operation ?? 'CONTAINS',
      };
    }),
  );

  useEffect(() => {
    setDisplayValue(value === '' ? '' : value);
  }, [value]);

  const onClick = (e: any) => {
    setCurrentContainerClassName(
      initialContainerClassName + selectedContainerClassName,
    );

    overlayRef?.current?.show(e, null);

    searchData();
  };

  const searchData = () => {
    if (!loadingOptions) {
      const wh = [] as any;

      filters
        .filter((f: any) => f.value)
        .forEach((f: any) => {
          wh.push({
            property: f.field,
            value: f.value,
            operation: f.operation,
          });
        });
      let where = undefined as unknown as any;
      switch (wh.length) {
        case 0: {
          where = undefined;
          break;
        }
        case 1: {
          where = {
            filter: wh[0],
          };
          break;
        }
        default: {
          where = {
            AND: [],
          };
          wh.forEach((f: any) => {
            where['AND'].push({
              filter: f,
            });
          });
        }
      }

      props.searchData(where);
    }
  };

  const onClear = () => {
    overlayRef?.current?.hide();
    setSelectedItemClassName(
      initialSelectedItemClassName + selectedItemEmptyClassName,
    );
    props.onItemSelected(undefined);
  };

  return (
    <div className={styles.container}>
      <label htmlFor='ownerFullName'>Owner</label>
      {triggerType === 'dropdown' && (
        <div
          ref={containerRef}
          className={classNames(
            currentContainerClassName,
            styles.searchModalTrigger,
            styles.xxxs,
          )}
          onClick={onClick}
        >
          <span ref={labelRef} className={selectedItemClassName}>
            {displayValue}
          </span>

          {displayValue !== '' && (
            <span
              className='flex align-items-center pl-2 pr-2'
              style={{ color: 'black' }}
              onClick={onClear}
            >
              {/*<FontAwesomeIcon icon={faTimes} />*/}
            </span>
          )}

          <div
            className='flex align-items-center p-dropdown-trigger'
            onClick={(e: any) => labelRef?.current?.click()}
          >
            <span className='p-dropdown-trigger-icon p-clickable pi pi-chevron-down'></span>
          </div>
        </div>
      )}

      {triggerType === 'button' && (
        <Button onClick={onClick} mode='primary'>
          <FontAwesomeIcon icon={props.buttonIcon} className='mr-2' />
          {props.buttonLabel}
        </Button>
      )}

      <OverlayPanel
        ref={overlayRef}
        style={{ width: overlayWidth }}
        onHide={() => setCurrentContainerClassName(initialContainerClassName)}
      >
        <div className={styles.searchModalWrapper}>
          {filters.length === 0 && (
            <div className='font-bold uppercase w-full mb-3'>
              {props.resourceLabel}
            </div>
          )}

          {filters.length > 0 && (
            <div>
              <div className='w-full mb-3'>
                Search{' '}
                <span className='font-bold lowercase'>
                  {props.resourceLabel}
                </span>{' '}
                by
              </div>

              <div className={styles.filtersContainer}>
                <div className={styles.filterOption}>
                  {filters?.map((f: any) => {
                    return (
                      <div className={styles.filterInput} key={f.field}>
                        <span className='flex flex-grow-1 mr-3'>
                          <Input
                            label={f.label}
                            inputSize='xxs'
                            className='w-full'
                            onChange={(e: any) => {
                              setFilters(
                                filters.map((fv: any) => {
                                  if (fv.field === f.field) {
                                    fv.value = e.target.value;
                                  }
                                  return fv;
                                }),
                              );
                            }}
                          />
                        </span>
                      </div>
                    );
                  })}
                </div>
                <div className='flex align-items-center'>
                  <Button mode='primary' onClick={searchData}>
                    Search
                  </Button>
                </div>
              </div>
            </div>
          )}

          {loadingOptions && (
            <div className='p-4'>
              <ul className='m-0 p-0'>
                <li className='mb-3'>
                  <div className='flex'>
                    <div style={{ flex: '1' }}>
                      <Skeleton width='100%' className='mb-2'></Skeleton>
                      <Skeleton width='75%'></Skeleton>
                    </div>
                  </div>
                </li>
                <li className=''>
                  <div className='flex'>
                    <div style={{ flex: '1' }}>
                      <Skeleton width='100%' className='mb-2'></Skeleton>
                      <Skeleton width='75%'></Skeleton>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
          )}

          {!loadingOptions &&
            options.length === 0 &&
            filters.filter((f: any) => f.value).length === 0 && (
              <div>There is no data to display</div>
            )}
          {!loadingOptions &&
            options.length === 0 &&
            filters.filter((f: any) => f.value).length > 0 && (
              <div>No data match your search criteria</div>
            )}

          <ul className={styles.optionsContainer}>
            {!loadingOptions &&
              options.map((c: any) => {
                return (
                  <li
                    role='button'
                    tabIndex={0}
                    key={c.id}
                    className={styles.option}
                    onClick={() => {
                      props.onItemSelected(c);
                      overlayRef?.current?.hide();
                      setSelectedItemClassName(initialSelectedItemClassName);
                      setFilters(
                        filters.map((fv: any) => {
                          fv.value = undefined;
                          return fv;
                        }),
                      );
                    }}
                  >
                    {props.itemTemplate(c)}
                  </li>
                );
              })}
          </ul>
        </div>

        {!loadingData && totalElements > maxResults && (
          <>
            <div>{totalElements} elements match your search term</div>
            <div>Improve your search term to narrow down the results</div>
          </>
        )}
        {/*<ChevronDown style={{ transform: 'scale(0.8)' }} />*/}
      </OverlayPanel>
    </div>
  );
};

import React from 'react';
import Image from 'next/image';
import styles from './settings-integration-item.module.scss';
import classNames from 'classnames';
import { Button, DebouncedInput } from '../../atoms';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'react-toastify';

interface FieldDefinition {
  name: string;
  label: string;
}

interface Props {
  icon: string;
  name: string;
  state: string;

  fields?: FieldDefinition[];

  onCancel?: () => void;
  settingsChanged?: () => void;
  onRevoke?: () => Promise<any>;
  onSave?: (data: any) => Promise<any>;
}

export const SettingsIntegrationItem = ({
  icon,
  name,
  state,
  fields,
  onCancel,
  settingsChanged,
  onRevoke,
  onSave,
}: Props) => {
  const [collapsed, setCollapsed] = React.useState(true);

  const { handleSubmit, setValue, getValues, control, reset } = useForm({
    defaultValues: fields?.map(({ name }) => {
      return { name: '' };
    }),
  });

  let stateFillColor = '';
  switch (state) {
    case 'ACTIVE':
      stateFillColor = 'green';
      break;
    case 'INACTIVE':
      stateFillColor = 'orange';
      break;
    case 'ERROR':
      stateFillColor = 'red';
      break;
  }

  const onRewokeSettings = () => {
    onRevoke &&
      onRevoke()
        .then(() => {
          setCollapsed(true);
          toast.success('Settings updated successfully!');
          settingsChanged && settingsChanged();
        })
        .catch(() => {
          toast.error(
            'There was a problem on our side and we are doing our best to solve it!',
          );
        });
  };

  const onSubmit = () => {
    onSave &&
      onSave(getValues())
        .then(() => {
          setCollapsed(true);
          toast.success('Settings updated successfully!');
          settingsChanged && settingsChanged();
        })
        .catch(() => {
          toast.error(
            'There was a problem on our side and we are doing our best to solve it!',
          );
        });
  };

  return (
    <div className={styles.settingsItem}>
      <div className={styles.settingsInfo}>
        <div className={styles.icon}>
          <Image alt='' src={icon} fill objectFit={'contain'} />
        </div>

        <div className={styles.name}>{name}</div>

        {/*TODO show state column all the time*/}
        {state === 'ACTIVE' && (
          <div className={styles.state}>
            <div className={styles.stateIcon}>
              <svg height='32' width='32'>
                <circle
                  cx='20'
                  cy='20'
                  r='12'
                  stroke='white'
                  stroke-width='3'
                  fill={stateFillColor}
                />
              </svg>
            </div>
            <div className={styles.stateText}>{state}</div>
          </div>
        )}

        <div className={styles.actions}>
          {collapsed && (
            <>
              {state === 'ACTIVE' && (
                <Button
                  onClick={() => {
                    setCollapsed(false);
                  }}
                  mode='secondary'
                >
                  Edit
                </Button>
              )}

              {state === 'INACTIVE' && (
                <Button
                  onClick={() => {
                    setCollapsed(false);
                  }}
                  mode='primary'
                >
                  Enable
                </Button>
              )}
            </>
          )}
          {!collapsed && (
            <>
              {state === 'ACTIVE' && (
                <>
                  <Button
                    onClick={() => {
                      setCollapsed(true);
                      onCancel && onCancel();
                      reset();
                    }}
                    mode='secondary'
                    style={{ marginRight: '10px' }}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={() => {
                      onRewokeSettings();
                    }}
                    mode='danger'
                    style={{ marginRight: '10px' }}
                  >
                    Revoke
                  </Button>
                  <Button
                    onClick={() => {
                      onSubmit();
                    }}
                    mode='primary'
                  >
                    Done
                  </Button>
                </>
              )}

              {state === 'INACTIVE' && (
                <>
                  <Button
                    onClick={() => {
                      setCollapsed(true);
                      onCancel && onCancel();
                      reset();
                    }}
                    mode='secondary'
                    style={{ marginRight: '10px' }}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={() => {
                      onSubmit();
                    }}
                    mode='primary'
                  >
                    Done
                  </Button>
                </>
              )}
            </>
          )}
        </div>
      </div>

      <div
        className={classNames(
          styles.settingsDetails,
          {
            [styles.collapsed]: collapsed,
          },
          {
            [styles.expanded]: !collapsed,
          },
        )}
      >
        {!collapsed && (
          <>
            <div className={styles.settingsDetailsContent}>
              {!fields && (
                <div
                  style={{ margin: '20px 0px 20px 60px', fontWeight: 'bold' }}
                >
                  Contact us!
                </div>
              )}

              {fields &&
                fields.map(
                  (fieldDefinition: FieldDefinition, index: number) => (
                    <div className={styles.field} key={fieldDefinition.name}>
                      <div className={styles.fieldLabel}>
                        {fieldDefinition.label}
                      </div>

                      <Controller
                        name={`${fieldDefinition.name}` as any}
                        control={control}
                        render={({ field }) => (
                          <input
                            value={
                              state === 'ACTIVE'
                                ? '******************'
                                : field.value as any
                            }
                            disabled={state === 'ACTIVE'}
                            className={styles.input}
                            onChange={({ target: { value } }) => {
                              field.onChange(value);
                            }}
                          />
                        )}
                      />
                    </div>
                  ),
                )}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

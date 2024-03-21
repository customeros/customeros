import React from 'react';
import Image from 'next/image';
import { useForm, Controller } from 'react-hook-form';

import { DeleteIntegrationSettings, UpdateIntegrationSettings } from 'services';

import { Fade } from '@ui/transitions/Fade';
import { Input } from '@ui/form/Input/Input2';
import { Button } from '@ui/form/Button/Button';
import { Collapse } from '@ui/transitions/Collapse';
import { Textarea } from '@ui/form/Textarea/Textarea';
import { SlideFade } from '@ui/transitions/SlideFade';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

interface FieldDefinition {
  name: string;
  label: string;
  textarea?: boolean;
}

interface Props {
  icon: string;
  name: string;
  state: string;
  identifier: string;
  onCancel?: () => void;
  onEnable?: () => void;
  onDisable?: () => void;
  fields?: FieldDefinition[];
  settingsChanged?: () => void;
}

export const SettingsIntegrationItem = ({
  icon,
  identifier,
  name,
  state,
  fields,
  onCancel,
  onEnable,
  onDisable,
  settingsChanged,
}: Props) => {
  const [collapsed, setCollapsed] = React.useState(true);

  const { getValues, control, reset } = useForm({
    defaultValues: fields?.map(({ name }) => {
      return { name: '' };
    }),
  });

  const onRevoke = () => {
    DeleteIntegrationSettings(identifier)
      .then(() => {
        setCollapsed(true);
        toastSuccess(
          'Settings updated successfully!',
          `${identifier}-integration-revoked`,
        );
        settingsChanged && settingsChanged();
      })
      .catch(() => {
        toastError(
          'There was a problem on our side and we are doing our best to solve it!',
          `${identifier}-integration-revoke-failed`,
        );
      });
  };

  const onSave = () => {
    UpdateIntegrationSettings(identifier, getValues())
      .then(() => {
        setCollapsed(true);
        toastSuccess(
          'Settings updated successfully!',
          `${identifier}-integration-settings-saved`,
        );
        settingsChanged && settingsChanged();
      })
      .catch(() => {
        toastError(
          'There was a problem on our side and we are doing our best to solve it!',
          `${identifier}-integration-settings-saved-failed`,
        );
      });
  };

  return (
    <div className='flex space-y-1 flex-col'>
      <div className='flex justify-between my-2'>
        <div className='flex items-center'>
          <Image className='mr-2' alt='' src={icon} width={20} height={20} />

          <span className='self-center text-md font-medium'>{name}</span>
        </div>

        <div>
          {collapsed && (
            <div className='flex space-x-1'>
              {state === 'ACTIVE' && (
                <Button
                  size='sm'
                  variant='outline'
                  colorScheme='gray'
                  onClick={() => {
                    if (onDisable) {
                      onDisable();
                    } else {
                      setCollapsed(false);
                    }
                  }}
                >
                  Edit
                </Button>
              )}

              {state === 'INACTIVE' && (
                <Button
                  size='sm'
                  variant='outline'
                  colorScheme='gray'
                  onClick={() => {
                    // If onEnable is present -> we're using the integration.app flows
                    if (onEnable) {
                      onEnable();
                    } else {
                      setCollapsed(false);
                    }
                  }}
                >
                  Enable
                </Button>
              )}
            </div>
          )}

          {!collapsed && (
            <Collapse in={!collapsed} style={{ overflow: 'unset' }}>
              <Fade in={!collapsed}>
                <div className='flex space-x-1'>
                  {state === 'ACTIVE' && (
                    <>
                      <Button
                        size='sm'
                        variant='outline'
                        colorScheme='gray'
                        onClick={() => {
                          setCollapsed(true);
                          onCancel && onCancel();
                          reset();
                        }}
                        style={{ marginRight: '10px' }}
                      >
                        Cancel
                      </Button>
                      <Button
                        size='sm'
                        variant='outline'
                        colorScheme='error'
                        onClick={onRevoke}
                        style={{ marginRight: '10px' }}
                      >
                        Revoke
                      </Button>
                      <Button
                        size='sm'
                        variant='outline'
                        colorScheme='success'
                        onClick={onSave}
                      >
                        Done
                      </Button>
                    </>
                  )}

                  {state === 'INACTIVE' && (
                    <>
                      <Button
                        size='sm'
                        variant='outline'
                        colorScheme='gray'
                        onClick={() => {
                          setCollapsed(true);
                          onCancel && onCancel();
                          reset();
                        }}
                        style={{ marginRight: '10px' }}
                      >
                        Cancel
                      </Button>
                      <Button
                        size='sm'
                        variant='outline'
                        onClick={onSave}
                        colorScheme='success'
                      >
                        Done
                      </Button>
                    </>
                  )}
                </div>
              </Fade>
            </Collapse>
          )}
        </div>
      </div>

      <div>
        <Collapse
          in={!collapsed}
          style={{ overflow: 'hidden' }}
          delay={{
            exit: 2,
          }}
        >
          <SlideFade in={!collapsed}>
            <>
              {!fields && (
                <span className='w-full m-5 mt-0 font-medium'>Contact us!</span>
              )}

              {fields &&
                fields.map((fieldDefinition: FieldDefinition) => (
                  <div
                    className=' flex mb-2 items-center'
                    key={fieldDefinition.name}
                  >
                    <label
                      className='mr-3 whitespace-nowrap'
                      htmlFor={fieldDefinition.name}
                    >
                      {fieldDefinition.label}
                    </label>

                    <Controller
                      // @ts-expect-error TODO: react-inverted-form should be used instead of hook-form
                      name={`${fieldDefinition.name}`}
                      control={control}
                      render={({ field }) => {
                        if (fieldDefinition.textarea) {
                          return (
                            <Textarea
                              id={fieldDefinition.name}
                              value={
                                state === 'ACTIVE'
                                  ? '******************'
                                  : (field.value as string)
                              }
                              disabled={state === 'ACTIVE'}
                              rows={1}
                              onChange={({ target: { value } }) => {
                                field.onChange(value);
                              }}
                              border
                            />
                          );
                        } else {
                          return (
                            <Input
                              id={fieldDefinition.name}
                              value={
                                state === 'ACTIVE'
                                  ? '******************'
                                  : (field.value as string)
                              }
                              disabled={state === 'ACTIVE'}
                              onChange={({ target: { value } }) => {
                                field.onChange(value);
                              }}
                              border
                            />
                          );
                        }
                      }}
                    />
                  </div>
                ))}
            </>
          </SlideFade>
        </Collapse>
      </div>
    </div>
  );
};

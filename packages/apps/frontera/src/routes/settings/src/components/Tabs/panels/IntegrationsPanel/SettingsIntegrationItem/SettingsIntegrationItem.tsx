import { useState } from 'react';
import { useForm, Controller } from 'react-hook-form';

import { Input } from '@ui/form/Input/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { AutoresizeTextarea } from '@ui/form/Textarea/AutoresizeTextarea';
import {
  CollapsibleRoot,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@ui/transitions/Collapse/Collapse';

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
  onSuccess?: () => void;
  fields?: FieldDefinition[];
  isIntegrationApp?: boolean;
}

export const SettingsIntegrationItem = ({
  icon,
  name,
  state,
  fields,
  onCancel,
  onEnable,
  onDisable,
  onSuccess,
  identifier,
  isIntegrationApp,
}: Props) => {
  const store = useStore();
  const [collapsed, setCollapsed] = useState(true);

  const { getValues, control, reset } = useForm({
    defaultValues: fields?.map(() => {
      return { name: '' };
    }),
  });

  const onRevoke = () => {
    store.settings.integrations.delete(identifier);
    onSuccess?.();
  };

  const onSave = () => {
    store.settings.integrations.update(identifier, getValues());
    onSuccess?.();
  };

  return (
    <CollapsibleRoot
      open={!collapsed}
      className='flex space-y-1 flex-col'
      onOpenChange={(value) => !isIntegrationApp && setCollapsed(!value)}
    >
      <div className='flex flex-row justify-between my-1'>
        <div className='flex items-center'>
          <img alt='' src={icon} width={20} height={20} className='mr-2' />
          <span className='self-center text-[14px]'>{name}</span>
        </div>
        <CollapsibleTrigger asChild={false} className='w-fit'>
          {collapsed && (
            <div className='flex space-x-1'>
              {state === 'ACTIVE' && collapsed && (
                <Button
                  size='xs'
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
                  size='xs'
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
            <div className='flex space-x-1'>
              {state === 'ACTIVE' && (
                <>
                  <Button
                    size='xs'
                    variant='outline'
                    colorScheme='gray'
                    style={{ marginRight: '10px' }}
                    onClick={() => {
                      setCollapsed(true);
                      onCancel && onCancel();
                      reset();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    size='xs'
                    variant='outline'
                    onClick={onRevoke}
                    colorScheme='error'
                    style={{ marginRight: '10px' }}
                  >
                    Revoke
                  </Button>
                  <Button
                    size='xs'
                    onClick={onSave}
                    variant='outline'
                    colorScheme='success'
                  >
                    Done
                  </Button>
                </>
              )}

              {state === 'INACTIVE' && (
                <>
                  <Button
                    size='xs'
                    variant='outline'
                    colorScheme='gray'
                    style={{ marginRight: '10px' }}
                    onClick={() => {
                      setCollapsed(true);
                      onCancel && onCancel();
                      reset();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    size='xs'
                    onClick={onSave}
                    variant='outline'
                    colorScheme='success'
                  >
                    Done
                  </Button>
                </>
              )}
            </div>
          )}
        </CollapsibleTrigger>
      </div>

      <CollapsibleContent>
        {!fields && (
          <span className='w-full m-5 mt-0 font-medium'>Contact us!</span>
        )}
        {fields &&
          fields.map((fieldDefinition: FieldDefinition) => (
            <div key={fieldDefinition.name} className='flex mb-2 items-center'>
              <label
                htmlFor={fieldDefinition.name}
                className='mr-3 whitespace-nowrap'
              >
                {fieldDefinition.label}
              </label>

              <Controller
                control={control}
                // @ts-expect-error TODO: react-inverted-form should be used instead of hook-form
                name={`${fieldDefinition.name}`}
                render={({ field }) => {
                  if (fieldDefinition.textarea) {
                    return (
                      <AutoresizeTextarea
                        rows={1}
                        id={fieldDefinition.name}
                        className='border-gray-200'
                        disabled={state === 'ACTIVE'}
                        onChange={({ target: { value } }) => {
                          field.onChange(value);
                        }}
                        value={
                          state === 'ACTIVE'
                            ? '******************'
                            : (field.value as string)
                        }
                      />
                    );
                  } else {
                    return (
                      <Input
                        id={fieldDefinition.name}
                        className='border-gray-200'
                        disabled={state === 'ACTIVE'}
                        onChange={({ target: { value } }) => {
                          field.onChange(value);
                        }}
                        value={
                          state === 'ACTIVE'
                            ? '******************'
                            : (field.value as string)
                        }
                      />
                    );
                  }
                }}
              />
            </div>
          ))}
      </CollapsibleContent>
    </CollapsibleRoot>
  );
};

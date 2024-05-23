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
      className='flex space-y-1 flex-col'
      open={!collapsed}
      onOpenChange={(value) => !isIntegrationApp && setCollapsed(!value)}
    >
      <div className='flex flex-row justify-between my-2'>
        <div className='flex items-center'>
          <img className='mr-2' alt='' src={icon} width={20} height={20} />
          <span className='self-center text-md font-medium'>{name}</span>
        </div>
        <CollapsibleTrigger className='w-fit' asChild={false}>
          {collapsed && (
            <div className='flex space-x-1'>
              {state === 'ACTIVE' && collapsed && (
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
          )}
        </CollapsibleTrigger>
      </div>

      <CollapsibleContent>
        {!fields && (
          <span className='w-full m-5 mt-0 font-medium'>Contact us!</span>
        )}
        {fields &&
          fields.map((fieldDefinition: FieldDefinition) => (
            <div className='flex mb-2 items-center' key={fieldDefinition.name}>
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
                      <AutoresizeTextarea
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
                        className='border-gray-200'
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
                        className='border-gray-200'
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

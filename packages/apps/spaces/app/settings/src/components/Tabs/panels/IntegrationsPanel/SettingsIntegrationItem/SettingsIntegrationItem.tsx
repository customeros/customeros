import React from 'react';
import { Button } from '@ui/form/Button';
import { Controller, useForm } from 'react-hook-form';
import { DeleteIntegrationSettings, UpdateIntegrationSettings } from 'services';
import {
  Box,
  Collapse,
  Fade,
  Flex,
  SlideFade,
  Text,
  Textarea,
} from '@chakra-ui/react';
import { Input } from '@ui/form/Input';
import { Image } from '@ui/media/Image';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

interface FieldDefinition {
  name: string;
  label: string;
  textarea?: boolean;
}

interface Props {
  icon: string;
  identifier: string;
  name: string;
  state: string;
  fields?: FieldDefinition[];
  onCancel?: () => void;
  onEnable?: () => void;
  onDisable?: () => void;
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
    <Flex direction='column'>
      <Flex justifyContent='space-between' my={2}>
        <Flex alignItems='center'>
          <Image alt='' src={icon} height={5} width={5} mr={2} />

          <Text fontWeight='medium' alignSelf='center' fontSize='md'>
            {name}
          </Text>
        </Flex>

        <Box>
          {collapsed && (
            <Flex>
              {state === 'ACTIVE' && (
                <Button
                  size='sm'
                  variant='outline'
                  onClick={() => {
                    if (onDisable) {
                      onDisable();
                    } else {
                      setCollapsed(false);
                    }
                  }}
                >
                  {onDisable ? 'Revoke' : 'Edit'}
                </Button>
              )}

              {state === 'INACTIVE' && (
                <Button
                  size='sm'
                  variant='outline'
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
            </Flex>
          )}

          {!collapsed && (
            <Collapse in={!collapsed} style={{ overflow: 'unset' }}>
              <Fade in={!collapsed}>
                <Flex>
                  {state === 'ACTIVE' && (
                    <>
                      <Button
                        size='sm'
                        variant='outline'
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
                        colorScheme='green'
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
                        colorScheme='green'
                      >
                        Done
                      </Button>
                    </>
                  )}
                </Flex>
              </Fade>
            </Collapse>
          )}
        </Box>
      </Flex>

      <Box>
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
                <Text m={5} mt={0} fontWeight='medium'>
                  Contact us!
                </Text>
              )}

              {fields &&
                fields.map((fieldDefinition: FieldDefinition) => (
                  <Flex key={fieldDefinition.name} alignItems='center' mb={2}>
                    <Text
                      whiteSpace='nowrap'
                      mr={3}
                      as='label'
                      htmlFor={fieldDefinition.name}
                    >
                      {fieldDefinition.label}
                    </Text>

                    <Controller
                      name={`${fieldDefinition.name}` as any}
                      control={control}
                      render={({ field }) => {
                        if (fieldDefinition.textarea) {
                          return (
                            <Textarea
                              id={fieldDefinition.name}
                              value={
                                state === 'ACTIVE'
                                  ? '******************'
                                  : (field.value as any)
                              }
                              disabled={state === 'ACTIVE'}
                              rows={1}
                              onChange={({ target: { value } }) => {
                                field.onChange(value);
                              }}
                              borderBottom='1px solid'
                              borderColor='gray.200'
                            />
                          );
                        } else {
                          return (
                            <Input
                              borderBottom='1px solid'
                              borderColor='gray.200'
                              id={fieldDefinition.name}
                              value={
                                state === 'ACTIVE'
                                  ? '******************'
                                  : (field.value as any)
                              }
                              disabled={state === 'ACTIVE'}
                              onChange={({ target: { value } }) => {
                                field.onChange(value);
                              }}
                            />
                          );
                        }
                      }}
                    />
                  </Flex>
                ))}
            </>
          </SlideFade>
        </Collapse>
      </Box>
    </Flex>
  );
};

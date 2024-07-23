import React, { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Switch } from '@ui/form/Switch';
import { Eye } from '@ui/media/icons/Eye.tsx';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { EyeOff } from '@ui/media/icons/EyeOff.tsx';
import { Button } from '@ui/form/Button/Button.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { InputGroup, RightElement } from '@ui/form/InputGroup';

export const LinkedInSettings = observer(() => {
  const store = useStore();
  const [isOpen, setIsOpen] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [validationError, setValidationError] = useState({
    linkedIn: false,
    password: false,
  });
  const [linkedIn, setLinkedIn] = useState('');
  const [password, setPassword] = useState('');

  const handleClose = () => {
    setIsOpen(false);
    setLinkedIn('');
    setPassword('');
    setValidationError({ linkedIn: false, password: false });
  };
  const handleSaveLinkedInCredentials = async () => {
    if (!linkedIn.length || !password.length) {
      setValidationError({
        linkedIn: !linkedIn.length,
        password: !password.length,
      });

      return;
    }

    store.settings.integrations.update('linkedin', {
      linkedInCredential: linkedIn,
      linkedInPassword: password,
    });
    handleClose();
  };

  const isActive =
    store.settings.integrations.value?.linkedin?.state === 'ACTIVE';

  return (
    <>
      <article className='flex-col flex relative max-w-[550px] px-6 '>
        <div className='flex items-center w-full'>
          <div className='flex items-center gap-1'>
            <h2 className='text-gray-700 text-sm font-medium'>LinkedIn</h2>
          </div>
          <div className='w-full border-b border-gray-100 mx-2' />

          <Tooltip label={isActive ? 'Remove LinkedIn account' : ''}>
            <div className='flex items-center'>
              <Switch
                disabled={store.settings.integrations.isBootstrapping}
                isChecked={isOpen || isActive}
                colorScheme='primary'
                size='sm'
                onChange={(isChecked) => {
                  if (isActive) {
                    store.settings.integrations.delete('linkedin');

                    return;
                  }

                  setIsOpen(isChecked);
                }}
              />
            </div>
          </Tooltip>
        </div>

        <p className='line-clamp-2 mt-2 mb-3 text-sm'>
          Import your LinkedIn connections by providing your email and password
        </p>

        {isOpen && (
          <>
            <label className='font-semibold text-sm mb-2'>
              Email or Phone
              <Input
                name='emailOrPhone'
                placeholder='olivia@untitledui.com'
                autoComplete='off'
                size='xs'
                className={cn(
                  'overflow-hidden overflow-ellipsis font-normal',
                  validationError.linkedIn && 'border-error-600',
                )}
                value={linkedIn}
                onChange={(e) => {
                  setLinkedIn(e.target.value);
                }}
              />
            </label>
            {validationError.linkedIn && (
              <p className='text-xs text-error-600 pt-1 -mt-1 mb-2'>
                Please enter an email or phone number
              </p>
            )}
            <label className='font-semibold text-sm mb-2 group'>
              Password
              <InputGroup
                className={cn(validationError.password && 'border-error-600')}
              >
                <Input
                  name='linkedInPassword'
                  placeholder='*********'
                  size='xs'
                  type={showPassword ? 'text' : 'password'}
                  className={
                    'overflow-hidden overflow-ellipsis font-normal border-none'
                  }
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <RightElement>
                  <IconButton
                    variant='ghost'
                    className='invisible group-hover:visible'
                    size='xs'
                    aria-label='Show password'
                    icon={showPassword ? <EyeOff /> : <Eye />}
                    onClick={() => setShowPassword(!showPassword)}
                  />
                </RightElement>
              </InputGroup>
            </label>
            {validationError.password && (
              <p className='text-xs text-error-600 pt-1 -mt-1'>
                Please enter a password
              </p>
            )}
            <div className='flex justify-end gap-2'>
              <Button size='xs' onClick={handleClose}>
                Cancel
              </Button>
              <Button
                size='xs'
                colorScheme='primary'
                onClick={handleSaveLinkedInCredentials}
              >
                Save
              </Button>
            </div>
          </>
        )}
      </article>
    </>
  );
});

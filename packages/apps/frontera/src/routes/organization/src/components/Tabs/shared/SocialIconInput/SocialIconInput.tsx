import React, { useRef, useMemo, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Social } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import {
  InputGroup,
  LeftElement,
  InputGroupProps,
} from '@ui/form/InputGroup/InputGroup';

import { SocialIcon } from './SocialIcons';
import { SocialInput } from './SocialInput';

interface FormSocialInputProps extends InputGroupProps {
  name: string;
  isReadOnly?: boolean;
  placeholder?: string;
  organizationId: string;
  leftElement?: React.ReactNode;
}

export const SocialIconInput = observer(
  ({
    name,
    leftElement,
    isReadOnly,
    organizationId,
    ...rest
  }: FormSocialInputProps) => {
    const store = useStore();
    const [socialIconValue, setSocialIconValue] = useState('');
    const organization = store.organizations.value.get(organizationId);
    const _leftElement = useMemo(() => leftElement, [leftElement]);
    const newInputRef = useRef<HTMLInputElement>(null);

    const focusNewInput = () => newInputRef.current?.focus();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;
      const value = e.target.value;

      if (organization) {
        organization.update((org) => {
          const idx = organization?.value.socialMedia.findIndex(
            (s) => s.id === id,
          );

          if (idx !== -1) {
            org.socialMedia[idx].url = value;
          }

          return org;
        });
      }
    };

    const handleBlur = (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;

      organization?.update((org) => {
        const idx = organization?.value.socialMedia.findIndex(
          (s) => s.id === id,
        );

        if (org.socialMedia[idx].url === '') {
          org.socialMedia.splice(idx, 1);
          focusNewInput();
        }

        return org;
      });
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;

      organization?.update((org) => {
        const idx = org.socialMedia.findIndex((s) => s.id === id);
        const social = org.socialMedia[idx];

        if (!social) return org;
        if (social.url === '') {
          org.socialMedia.splice(idx, 1);
          focusNewInput();
        }

        return org;
      });
    };

    const handleNewSocial = () => {
      const value = newInputRef.current?.value;
      if (!value) return;
      organization?.update((org) => {
        org.socialMedia.push({
          id: crypto.randomUUID(),
          url: value,
        } as Social);

        return org;
      });
      newInputRef.current!.value = '';
      setSocialIconValue('');
    };

    return (
      <>
        {organization?.value.socialMedia?.map(({ id, url }) => (
          <SocialInput
            name='socialMedia'
            id={id}
            key={id}
            value={url}
            onBlur={handleBlur}
            onChange={handleChange}
            isReadOnly={isReadOnly}
            onKeyDown={handleKeyDown}
            leftElement={_leftElement}
          />
        ))}

        {!isReadOnly && (
          <InputGroup {...rest}>
            {leftElement && (
              <LeftElement>
                <SocialIcon url={socialIconValue}>{leftElement}</SocialIcon>
              </LeftElement>
            )}
            <Input
              name='socialMedia'
              className='border-b border-transparent hover:border-transparent hover:border-b-none text-md focus:hover:border-b focus:hover:border-transparent focus:border-b focus:border-transparent'
              ref={newInputRef}
              onBlur={handleNewSocial}
              onChange={(e) => setSocialIconValue(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleNewSocial();
                }
              }}
              {...rest}
            />
          </InputGroup>
        )}
      </>
    );
  },
);

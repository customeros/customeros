import { useParams } from 'react-router-dom';
import React, { useRef, useMemo, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { AddSocialLinkCase } from '@domain/usecases/organization-about-tab/add-social-link.usecase.ts';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

import { SocialIcon } from './SocialIcons';
import { SocialInput } from './SocialInput';

interface SocialIconInputProps {
  name?: string;
  dataTest?: string;
  isReadOnly?: boolean;
  placeholder?: string;
  leftElement?: React.ReactNode;
  value?: { label: string; value: string }[];
}
const addSocialLinkUseCase = new AddSocialLinkCase();

export const SocialIconInput = observer(
  ({
    value,
    name = 'socialMedia',
    leftElement,
    isReadOnly,
    dataTest,
    ...rest
  }: SocialIconInputProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.getById(id);

    const _leftElement = useMemo(() => leftElement, [leftElement]);
    const newInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
      if (organization) {
        addSocialLinkUseCase.setEntity(organization);
      }
    }, [organization?.id]);

    if (!organization || !organization?.value) return null;

    const handleSocialChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;
      const value = e.target.value;

      if (organization) {
        const idx = organization.value?.socialMedia.findIndex(
          (s) => s.id === id,
        );

        if (typeof idx !== 'number' || idx < 0) return;

        organization.draft();
        organization.value!.socialMedia[idx].url = value;
      }
    };

    const handleSocialBlur = (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;

      const idx = organization.value?.socialMedia.findIndex((s) => s.id === id);

      if (typeof idx !== 'number' || idx < 0) return;

      if (e.target?.value === '') {
        organization.draft();
        organization.value.socialMedia.splice(idx, 1);
        newInputRef.current?.focus();
        organization.commit();
      } else {
        organization.draft();
        organization.value!.socialMedia[idx].url = e.target.value;
        organization.commit();
      }
    };

    const handleSocialKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;

      const idx = organization.value?.socialMedia.findIndex((s) => s.id === id);

      if (typeof idx !== 'number' || idx < 0) return;
      const social = organization.value?.socialMedia[idx];

      if (!social) return;

      if ((e.target as HTMLInputElement).value === '') {
        organization.draft();
        organization.value.socialMedia.splice(idx, 1);
        organization.commit();
      } else {
        organization.draft();
        organization.value!.socialMedia[idx].url = (
          e.target as HTMLInputElement
        ).value;

        organization.commit();
      }
      newInputRef.current?.focus();
    };

    return (
      <>
        {value?.map(({ value: v, label: l }) => (
          <SocialInput
            id={v}
            key={v}
            value={l}
            name={name}
            dataTest={dataTest}
            isReadOnly={isReadOnly}
            onBlur={handleSocialBlur}
            leftElement={_leftElement}
            onChange={handleSocialChange}
            onKeyDown={handleSocialKeyDown}
            {...rest}
          />
        ))}

        {!isReadOnly && (
          <InputGroup>
            {leftElement && (
              <LeftElement>
                <SocialIcon url={addSocialLinkUseCase.url}>
                  {leftElement}
                </SocialIcon>
              </LeftElement>
            )}
            <Input
              name={name}
              ref={newInputRef}
              dataTest={dataTest}
              onBlur={addSocialLinkUseCase.submit}
              onChange={(e) => {
                addSocialLinkUseCase.setInputValue(e.target.value);
              }}
              onKeyDown={(e) => {
                e.stopPropagation();

                if (e.key === 'Enter') {
                  addSocialLinkUseCase.submit();
                }
              }}
              className='border-b border-transparent hover:border-transparent hover:border-b-none text-md focus:hover:border-b focus:hover:border-transparent focus:border-b focus:border-transparent'
              {...rest}
            />
          </InputGroup>
        )}
      </>
    );
  },
);

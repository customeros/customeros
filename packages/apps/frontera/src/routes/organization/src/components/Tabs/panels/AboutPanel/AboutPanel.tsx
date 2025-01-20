import { useRef } from 'react';
import { useParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding';
import { Users02 } from '@ui/media/icons/Users02';
import { Share07 } from '@ui/media/icons/Share07';
import { Target05 } from '@ui/media/icons/Target05';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building07 } from '@ui/media/icons/Building07';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { TruncatedText } from '@ui/presentation/TruncatedText/TruncatedText.tsx';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import {
  EntityType,
  OrganizationStage,
  OrganizationRelationship,
} from '@graphql/types';

import { Tags } from './components/tags';
import { Domains } from './components/Domains';
import { SocialMediaList } from '../../shared';
import { OwnerInput } from './components/owner';
import { Branches } from './components/branches';
import { stageOptions, getStageOptions, relationshipOptions } from './util';

const iconMap = {
  Customer: <ActivityHeart className='text-gray-500' />,
  Prospect: <Seeding className='text-gray-500' />,
  'Not a fit': <MessageXCircle className='text-gray-500' />,
  'Former Customer': <BrokenHeart className='text-gray-500' />,
  unknown: <AlignHorizontalCentre02 className='text-gray-500' />,
};

export const AboutPanel = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const [_, copyToClipboard] = useCopyToClipboard();
  const nameRef = useRef<HTMLInputElement | null>(null);

  const showParentRelationshipSelector = useFeatureIsOn(
    'show-parent-relationship-selector',
  );
  const parentRelationshipReadOnly = useFeatureIsOn(
    'parent-relationship-selector-read-only',
  );
  const orgNameReadOnly = useFeatureIsOn('org-name-readonly');

  const organization = store.organizations.getById(id);

  if (!organization || !organization?.value) return null;

  const selectedRelationshipOption = relationshipOptions.find(
    (option) => option.value === organization.value?.relationship,
  );

  const selectedStageOption = stageOptions.find(
    (option) => option.value === organization.value?.stage,
  );

  const applicableStageOptions = getStageOptions(
    organization.value?.relationship,
  );

  const handleCreateOption = (value: string) => {
    store.tags?.create(
      { name: value },
      {
        onSucces: (id) => {
          organization.draft();
          organization?.value?.tags?.push({
            name: value,
            colorCode: store.tags.getById(id)?.value?.colorCode ?? 'grayModern',
            metadata: {
              id,
            },
            entityType: EntityType.Organization,
          });
          organization.commit();
        },
      },
    );
  };

  const isEnriching = organization.isEnriching;

  return (
    <div className='flex pt-[6px] px-6 w-full h-full overflow-y-auto flex-1 bg-gray-25 rounded-2xl'>
      <div className='flex h-full flex-col  overflow-visible w-full'>
        {isEnriching && (
          <div className='flex items-center justify-start gap-2 border-[1px] text-sm border-grayModern-100 bg-grayModern-50 rounded-[4px] py-1 px-2 '>
            <Spinner
              label='enriching org'
              className='text-grayModern-300 fill-grayModern-500 size-4'
            />
            <span className='font-medium'>
              We're enriching this organizations's details...
            </span>
          </div>
        )}

        <div className='flex items-center justify-between'>
          <Input
            size='xs'
            name='name'
            ref={nameRef}
            autoComplete='off'
            variant='unstyled'
            dataTest='org-about-name'
            placeholder='Company name'
            disabled={orgNameReadOnly}
            onFocus={(e) => e.target.select()}
            value={organization?.value.name || ''}
            onChange={(e) => {
              organization.value.name = e.target.value;
            }}
            className='font-semibold text-[16px] mt-0.5 border-none overflow-hidden overflow-ellipsis'
            onBlur={() => {
              organization.draft();
              organization.commit();
            }}
          />
          {organization.value?.referenceId && (
            <div className='h-full ml-4'>
              <Tooltip asChild={false} label={'Copy ID'}>
                <Tag
                  variant='outline'
                  colorScheme='gray'
                  className='rounded-full cursor-pointer'
                  onClick={() => {
                    copyToClipboard(
                      organization.value?.referenceId ?? '',
                      'Reference ID copied ',
                    );
                  }}
                >
                  <TagLabel>{organization.value?.referenceId}</TagLabel>
                </Tag>
              </Tooltip>
            </div>
          )}
        </div>

        <Domains />

        <div className='flex flex-col w-full flex-1 items-start justify-start gap-3 mt-2'>
          {!!organization?.value?.description && (
            <TruncatedText
              maxLines={4}
              className={'text-sm'}
              data-test='org-about-description'
              text={organization.value.description}
            />
          )}

          <Tags
            dataTest='org-about-tags'
            inputPlaceholder='Search...'
            onCreate={handleCreateOption}
            placeholder='Organization tags'
            leftAccessory={<Tag01 className='mr-3 text-gray-500' />}
            value={
              organization.value.tags?.map((t) => ({
                value: t.metadata.id,
                label: t.name,
              })) ?? []
            }
            options={store.tags
              .getByEntityType(EntityType.Organization)
              .map((t) => ({
                value: t.id,
                label: t.value?.name,
              }))}
            onChange={(selection) => {
              const tags = selection
                .map((o) => store.tags.getById(o.value)?.value)
                .filter(Boolean);

              organization.draft();
              organization.value.tags = tags as TagDatum[];
              organization.commit();
            }}
          />

          <div className='flex items-center justify-center w-full '>
            <div
              data-test='org-about-relationship'
              className='flex-2 flex items-center'
            >
              <Menu>
                <Tooltip align='start' label='Relationship'>
                  <MenuButton
                    data-test='org-about-relationship'
                    className='min-h-[20px] text-md outline-none focus:outline-none items-center'
                  >
                    {
                      iconMap[
                        (selectedRelationshipOption?.label ??
                          'unknown') as keyof typeof iconMap
                      ]
                    }
                    {''}
                    <span
                      className={cn(
                        'ml-3 text-sm',
                        !selectedRelationshipOption?.label && 'text-gray-400',
                      )}
                    >
                      {selectedRelationshipOption?.label ?? 'Relationship'}
                    </span>
                  </MenuButton>
                </Tooltip>
                <MenuList side='bottom' align='start'>
                  {relationshipOptions.map((option) => (
                    <MenuItem
                      key={option.value}
                      onClick={() => {
                        organization.value!.relationship = option.value;
                        organization.value!.stage = match(option.value)
                          .with(
                            OrganizationRelationship.Prospect,
                            () => OrganizationStage.Lead,
                          )
                          .with(
                            OrganizationRelationship.Customer,
                            () => OrganizationStage.InitialValue,
                          )
                          .with(
                            OrganizationRelationship.NotAFit,
                            () => OrganizationStage.Unqualified,
                          )
                          .with(
                            OrganizationRelationship.FormerCustomer,
                            () => OrganizationStage.Target,
                          )
                          .otherwise(() => undefined);

                        organization.commit();
                      }}
                    >
                      {iconMap[option.label as keyof typeof iconMap]}
                      {option.label}
                    </MenuItem>
                  ))}
                </MenuList>
              </Menu>
            </div>
            {selectedRelationshipOption?.value !==
              OrganizationRelationship.Customer && (
              <div
                data-test='org-about-stage'
                className='flex-1 flex items-center'
              >
                <Menu>
                  <Tooltip label='Stage' align='start'>
                    <MenuButton className='min-h-[20px] outline-none focus:outline-none'>
                      <Target05 className='text-gray-500 mb-0.5' />
                      <span className='ml-3 text-sm'>
                        {selectedStageOption?.label || 'Stage'}
                      </span>
                    </MenuButton>
                  </Tooltip>
                  <MenuList side='bottom' align='start'>
                    {applicableStageOptions.map((option) => (
                      <MenuItem
                        key={option.value}
                        onClick={() => {
                          organization.value!.stage = option.value;
                          organization.commit();
                        }}
                      >
                        {iconMap[option.label as keyof typeof iconMap]}
                        {option.label}
                      </MenuItem>
                    ))}
                  </MenuList>
                </Menu>
              </div>
            )}
          </div>
          <Tooltip align='start' label='Industry'>
            <p className='text-sm flex items-center cursor-default'>
              <Building07 className='text-gray-500 mr-3 ' />
              {organization?.value?.industryName ? (
                <span>{organization.value.industryName}</span>
              ) : (
                <span
                  className={'text-gray-400'}
                  data-test='org-about-industry'
                >
                  Industry not found yet
                </span>
              )}
            </p>
          </Tooltip>
          {organization.country && (
            <Tooltip align='start' label='Country'>
              <p className='text-sm flex items-center cursor-default'>
                <span className='flex items-center mr-3'>
                  {organization.value.locations?.[0]?.countryCodeA2 &&
                    flags[organization.value.locations?.[0]?.countryCodeA2]}
                </span>

                {organization.country}
              </p>
            </Tooltip>
          )}

          {typeof organization.value!.employees === 'number' && (
            <Tooltip align='start' label='Number of employees'>
              <p className='text-sm flex items-center cursor-default '>
                <Users02 className='text-gray-500 mr-3' />
                {organization.value!.employees} employees
              </p>
            </Tooltip>
          )}
          <OwnerInput
            id={id}
            dataTest='org-about-org-owner'
            owner={organization?.value.owner}
          />
          <SocialMediaList
            dataTest='org-about-social-link'
            leftElement={<Share07 className='text-gray-500' />}
            value={organization?.value.socialMedia.map((s) => ({
              value: s.id,
              label: s.url,
            }))}
          />

          {showParentRelationshipSelector &&
            organization?.value?.subsidiaries?.length > 0 && (
              <Branches id={id} isReadOnly={parentRelationshipReadOnly} />
            )}
        </div>
        {organization?.value.customerOsId && (
          <Tooltip label='Copy ID'>
            <span
              className='py-3 w-fit text-gray-400 cursor-pointer text-sm'
              onClick={() =>
                copyToClipboard(
                  organization.value?.customerOsId ?? '',
                  'CustomerOS ID copied',
                )
              }
            >
              CustomerOS ID: {organization?.value.customerOsId}
            </span>
          </Tooltip>
        )}
      </div>
    </div>
  );
});

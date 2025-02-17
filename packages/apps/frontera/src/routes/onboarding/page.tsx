import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { match } from 'ts-pattern';
import { useUnmount } from 'usehooks-ts';
import { observer } from 'mobx-react-lite';

import { Send03 } from '@ui/media/icons/Send03';
import { useStore } from '@shared/hooks/useStore';
import { Inbox01 } from '@ui/media/icons/Inbox01';
import { Shuffle01 } from '@ui/media/icons/Shuffle01';
import { Building07 } from '@ui/media/icons/Building07';
import { Settings01 } from '@ui/media/icons/Settings01';
import { WelcomePage } from '@ui/media/logos/WelcomePage';
import { PreviewCard } from '@shared/components/PreviewCard';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';
import { ShortcutsPanel } from '@shared/components/PreviewCard/components/ShortcutsPanel';

import { OnboardingCard } from './src/components/OnboardingCard';
import { DecorativePattern } from './src/assets/DecorativePattern';

export const OnboardingPage = observer(() => {
  const { globalCache, tableViewDefs, mailboxes, ui } = useStore();
  const navigate = useNavigate();

  const onboarding = globalCache.value?.user?.onboarding;

  useEffect(() => {
    if (
      !globalCache.isLoading &&
      !globalCache.value?.user.onboarding.showOnboardingPage
    ) {
      navigate('/finder');
    }
  }, [globalCache.value?.user.onboarding, globalCache.isLoading]);

  const handleOnboardingStepClick = (
    step: 'outbound' | 'crm' | 'mailstack',
  ) => {
    match(step)
      .with('outbound', () => {
        const flowsPreset = tableViewDefs.flowsPreset;

        globalCache.updateOnboardingDetails({
          onboardingOutboundStepCompleted: true,
        });
        navigate(`/finder?preset=${flowsPreset}`, {
          state: { fromOnboarding: true },
        });
      })
      .with('crm', () => {
        const contactsPreset = tableViewDefs.contactsPreset;

        globalCache.updateOnboardingDetails({
          onboardingCrmStepCompleted: true,
        });
        navigate(`/finder?preset=${contactsPreset}`, {
          state: { fromOnboarding: true },
        });
      })
      .with('mailstack', () => {
        const hasMailboxes = mailboxes.value.size > 0;

        globalCache.updateOnboardingDetails({
          onboardingMailstackStepCompleted: true,
        });

        if (hasMailboxes) {
          navigate('/settings?tab=mailboxes&view=buy');
        } else {
          navigate('/settings?tab=mailboxes');
        }
      })
      .otherwise(() => null);
  };

  useUnmount(() => {
    ui.setShortcutsPanel(false);
  });

  return (
    <div className='flex'>
      <div className='flex h-full flex-col w-full relative'>
        <div className='flex flex-col h-[323px] w-[500px] items-center self-center'>
          <div className='relative flex-col'>
            <WelcomePage className='w-[152px] h-[120px] absolute top-[20%] right-[34%]' />
            <DecorativePattern
              width={500}
              height={480}
              className='-mt-[100px]'
            />
          </div>
          <article className='flex flex-col items-center top-[-180px] relative'>
            <h1 className='font-bold text-xl'>Welcome!</h1>
            <p className='mt-2'>Let’s get you started.</p>
            <p className='mt-4'>Here are a few short paths to success:</p>
          </article>
        </div>

        <div className='flex gap-4 justify-center items-center '>
          <div>
            <OnboardingCard
              title='Outbound'
              timeInMinutes='2'
              isFeatured={false}
              contentTitle='Explore a sequence'
              subtitle='Automated, surgical, reputable'
              titleIcon={<Send03 className='text-inherit' />}
              contentIcon={<Shuffle01 className='text-inherit' />}
              onClick={() => handleOnboardingStepClick('outbound')}
              isCompleted={onboarding?.onboardingOutboundStepCompleted}
              description='Discover how to orchestrate a multichannel outbound campaign.'
            />{' '}
          </div>
          <div>
            <OnboardingCard
              title='CRM'
              isFeatured={true}
              timeInMinutes='3'
              contentTitle='Import a contact from LinkedIn'
              subtitle='All the customer data, for everyone'
              onClick={() => handleOnboardingStepClick('crm')}
              isCompleted={onboarding?.onboardingCrmStepCompleted}
              contentIcon={<LinkedinOutline className='size-5' />}
              titleIcon={<Building07 className='size-6 text-inherit' />}
              description='...then enrich their email and mobile number in one click'
            />{' '}
          </div>
          <div>
            <OnboardingCard
              title='Mailstack'
              timeInMinutes='5'
              isFeatured={false}
              titleIcon={<Inbox01 />}
              contentIcon={<Settings01 />}
              contentTitle='Setup my mailboxes'
              subtitle='High-deliverability email infrastructure'
              onClick={() => handleOnboardingStepClick('mailstack')}
              isCompleted={onboarding?.onboardingMailstackStepCompleted}
              description='Setup your high-reputation domains and senders, for pain-free outbound with optimal deliverability.'
            />
          </div>
        </div>
      </div>
      <div className='h-full'>
        {ui.showShortcutsPanel && (
          <PreviewCard>
            <ShortcutsPanel />
          </PreviewCard>
        )}
      </div>
    </div>
  );
});

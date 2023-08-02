import React, { useEffect, useRef } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { Button } from '@ui/form/Button';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { ComposeEmail } from '@organization/components/Timeline/events/email/compose-email/ComposeEmail';
import Envelope from '@spaces/atoms/icons/Envelope';

interface TimelineActionsProps {
  onScrollBottom: () => void;
}

export const TimelineActions: React.FC<TimelineActionsProps> = ({ onScrollBottom }) => {
  const [show, setShow] = React.useState(false);
  const virtuoso = useRef(null);

  useEffect(() => {
    if (show) {
      onScrollBottom();
    }
  }, [show]);
  const handleToggle = () => setShow(!show);
  return (
    <Box>
      <ButtonGroup
        mt={6}
        position='sticky'
        py={2}
        border='1px dashed var(--gray-200, #EAECF0)'
        p={2}
        borderRadius={30}
        bg='white'
        top='0'
        left={6}
        zIndex={1}
      >
        <Button
          variant='outline'
          onClick={() => handleToggle()}
          borderRadius='3xl'
          size='xs'
          leftIcon={<Envelope color='inherit' height={16} width={16} />}
        >
          Email
        </Button>
      </ButtonGroup>
      <Box
        bg={'#F9F9FB'}
        borderTop='1px dashed var(--gray-200, #EAECF0)'
        pt={6}
        pb={show ? 2 : 4}
        mt={-4}
      >
        {show && (
          <SlideFade in={true}>
            <Box
              ref={virtuoso}
              borderRadius={'md'}
              boxShadow={'lg'}
              m={6}
              mt={0}
              pb={4}
              bg={'white'}
              border='1px solid var(--gray-100, #F2F4F7)'
            >
              <ComposeEmail
                modal={false}
                to={[]}
                cc={[]}
                bcc={[]}
                from={[]}
                subject={''}
                emailContent={''}
              />
            </Box>
          </SlideFade>
        )}
      </Box>
    </Box>
  );
};

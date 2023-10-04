import { Center } from '@ui/layout/Center';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';

interface EmptyStateProps {
  onClick: () => void;
}

const EmptyState = ({ onClick }: EmptyStateProps) => {
  return (
    <Center
      h='100%'
      bg='white'
      borderRadius='2xl'
      border='1px solid'
      borderColor='gray.200'
    >
      <Flex position='absolute' top='15%'>
        <svg
          width='512'
          height='480'
          viewBox='0 0 512 480'
          fill='none'
          xmlns='http://www.w3.org/2000/svg'
        >
          <mask
            id='mask0_7_662'
            style={{ maskType: 'alpha' }}
            maskUnits='userSpaceOnUse'
            x='16'
            y='0'
            width='480'
            height='480'
          >
            <rect
              width='480'
              height='480'
              transform='translate(16)'
              fill='url(#paint0_radial_7_662)'
            />
          </mask>
          <g mask='url(#mask0_7_662)'>
            <circle cx='256' cy='240' r='47.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='79.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='111.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='143.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='143.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='175.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='207.5' stroke='#EAECF0' />
            <circle cx='256' cy='240' r='239.5' stroke='#EAECF0' />
          </g>
          <circle cx='256' cy='240' r='52' fill='#E9D7FE' />
          <g filter='url(#filter0_dd_7_662)'>
            <path
              fill-rule='evenodd'
              clip-rule='evenodd'
              d='M257.6 204C246.827 204 237.298 209.323 231.499 217.483C229.605 217.036 227.63 216.8 225.6 216.8C211.462 216.8 200 228.262 200 242.4C200 256.538 211.462 268 225.6 268L225.62 268H289.6C301.971 268 312 257.971 312 245.6C312 233.229 301.971 223.2 289.6 223.2C288.721 223.2 287.854 223.251 287.002 223.349C282.098 211.968 270.78 204 257.6 204Z'
              fill='#F9F5FF'
            />
            <ellipse
              cx='225.6'
              cy='242.4'
              rx='25.6'
              ry='25.6'
              fill='url(#paint1_linear_7_662)'
            />
            <circle
              cx='257.6'
              cy='236'
              r='32'
              fill='url(#paint2_linear_7_662)'
            />
            <ellipse
              cx='289.6'
              cy='245.6'
              rx='22.4'
              ry='22.4'
              fill='url(#paint3_linear_7_662)'
            />
          </g>
          <circle cx='201' cy='207' r='5' fill='#F4EBFF' />
          <circle cx='198' cy='297' r='7' fill='#F4EBFF' />
          <circle cx='325' cy='223' r='7' fill='#F4EBFF' />
          <circle cx='314' cy='196' r='4' fill='#F4EBFF' />
          <g filter='url(#filter1_b_7_662)'>
            <rect
              x='232'
              y='250'
              width='48'
              height='48'
              rx='24'
              fill='#6941C6'
              fill-opacity='0.4'
            />
            <path
              d='M248 278.242C246.794 277.435 246 276.06 246 274.5C246 272.156 247.792 270.231 250.08 270.019C250.548 267.172 253.02 265 256 265C258.98 265 261.452 267.172 261.92 270.019C264.208 270.231 266 272.156 266 274.5C266 276.06 265.206 277.435 264 278.242M252 278L256 274M256 274L260 278M256 274V283'
              stroke='white'
              stroke-width='2'
              stroke-linecap='round'
              stroke-linejoin='round'
            />
          </g>
          <defs>
            <filter
              id='filter0_dd_7_662'
              x='180'
              y='204'
              width='152'
              height='104'
              filterUnits='userSpaceOnUse'
              color-interpolation-filters='sRGB'
            >
              <feFlood flood-opacity='0' result='BackgroundImageFix' />
              <feColorMatrix
                in='SourceAlpha'
                type='matrix'
                values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
                result='hardAlpha'
              />
              <feMorphology
                radius='4'
                operator='erode'
                in='SourceAlpha'
                result='effect1_dropShadow_7_662'
              />
              <feOffset dy='8' />
              <feGaussianBlur stdDeviation='4' />
              <feColorMatrix
                type='matrix'
                values='0 0 0 0 0.0627451 0 0 0 0 0.0941176 0 0 0 0 0.156863 0 0 0 0.03 0'
              />
              <feBlend
                mode='normal'
                in2='BackgroundImageFix'
                result='effect1_dropShadow_7_662'
              />
              <feColorMatrix
                in='SourceAlpha'
                type='matrix'
                values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
                result='hardAlpha'
              />
              <feMorphology
                radius='4'
                operator='erode'
                in='SourceAlpha'
                result='effect2_dropShadow_7_662'
              />
              <feOffset dy='20' />
              <feGaussianBlur stdDeviation='12' />
              <feColorMatrix
                type='matrix'
                values='0 0 0 0 0.0627451 0 0 0 0 0.0941176 0 0 0 0 0.156863 0 0 0 0.08 0'
              />
              <feBlend
                mode='normal'
                in2='effect1_dropShadow_7_662'
                result='effect2_dropShadow_7_662'
              />
              <feBlend
                mode='normal'
                in='SourceGraphic'
                in2='effect2_dropShadow_7_662'
                result='shape'
              />
            </filter>
            <filter
              id='filter1_b_7_662'
              x='224'
              y='242'
              width='64'
              height='64'
              filterUnits='userSpaceOnUse'
              color-interpolation-filters='sRGB'
            >
              <feFlood flood-opacity='0' result='BackgroundImageFix' />
              <feGaussianBlur in='BackgroundImageFix' stdDeviation='4' />
              <feComposite
                in2='SourceAlpha'
                operator='in'
                result='effect1_backgroundBlur_7_662'
              />
              <feBlend
                mode='normal'
                in='SourceGraphic'
                in2='effect1_backgroundBlur_7_662'
                result='shape'
              />
            </filter>
            <radialGradient
              id='paint0_radial_7_662'
              cx='0'
              cy='0'
              r='1'
              gradientUnits='userSpaceOnUse'
              gradientTransform='translate(240 480) rotate(-90) scale(480 250.485)'
            >
              <stop />
              <stop offset='0.958333' stop-opacity='0' />
            </radialGradient>
            <linearGradient
              id='paint1_linear_7_662'
              x1='205.943'
              y1='225.486'
              x2='251.2'
              y2='268'
              gradientUnits='userSpaceOnUse'
            >
              <stop stop-color='#E9D7FE' />
              <stop offset='0.350715' stop-color='white' stop-opacity='0' />
            </linearGradient>
            <linearGradient
              id='paint2_linear_7_662'
              x1='233.029'
              y1='214.857'
              x2='289.6'
              y2='268'
              gradientUnits='userSpaceOnUse'
            >
              <stop stop-color='#E9D7FE' />
              <stop offset='0.350715' stop-color='white' stop-opacity='0' />
            </linearGradient>
            <linearGradient
              id='paint3_linear_7_662'
              x1='272.4'
              y1='230.8'
              x2='312'
              y2='268'
              gradientUnits='userSpaceOnUse'
            >
              <stop stop-color='#E9D7FE' />
              <stop offset='0.350715' stop-color='white' stop-opacity='0' />
            </linearGradient>
          </defs>
        </svg>
      </Flex>
      <Flex
        position='relative'
        flexDir='column'
        textAlign='center'
        top='-5%'
        align='center'
      >
        <Text
          color='gray.900'
          fontSize='md'
          fontWeight='semibold'
        >{`Let's get started`}</Text>
        <Text maxW='400px' fontSize='sm' color='gray.600'>
          Start seeing your customer conversations all in one place by adding an
          organization
        </Text>

        <Button onClick={onClick} mt='6' w='min-content' variant='outline'>
          Add Organization
        </Button>
      </Flex>
    </Center>
  );
};

export default EmptyState;

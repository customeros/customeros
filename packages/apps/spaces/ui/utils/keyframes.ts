import { keyframes } from '@chakra-ui/react';

const pulseOpacity = keyframes`
  from { opacity: 0.3; }
  to { opacity: 0.7; }
`;
const wave = keyframes`
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(3600deg);
  }
`;

export { pulseOpacity, wave };

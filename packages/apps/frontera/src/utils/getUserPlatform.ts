export const isUserPlatformMac = (): boolean => {
  if (navigator.userAgent) {
    return navigator.userAgent.toLowerCase().includes('mac');
  }

  return navigator.platform.toLowerCase().includes('mac');
};

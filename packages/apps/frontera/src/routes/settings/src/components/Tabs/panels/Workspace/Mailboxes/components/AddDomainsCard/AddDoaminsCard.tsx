import { useState } from 'react';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01';
import { ShoppingCartAdd } from '@ui/media/icons/ShopingCartAdd';
import {
  Card,
  CardFooter,
  CardHeader,
  CardContent,
} from '@ui/presentation/Card/Card';

export const AddDomainsCard = () => {
  const [brandName, setBrandName] = useState('');
  const [domainVariations, setDomainVariations] = useState<string[]>([]);
  const [isHovered, setIsHovered] = useState<number | null>(null);
  const [showSecondHalf, setShowSecondHalf] = useState(false);

  function generateDomainVariations(baseName: string) {
    if (!baseName) return [];
    const variations = [
      `get${baseName}`,
      `${baseName}hq`,
      `${baseName}shop`,
      `try${baseName}`,
      `${baseName}app`,
      `join${baseName}`,
      `${baseName}online`,
      `${baseName}solutions`,
      `${baseName}platform`,
      `${baseName}labs`,
      `${baseName}store`,
      `my${baseName}`,
      `${baseName}tech`,
      `go${baseName}`,
      `${baseName}systems`,
      `${baseName}now`,
      `${baseName}digital`,
      `${baseName}inc`,
      `${baseName}software`,
      `${baseName}cloud`,
    ];

    return variations.map((variation) => `${variation}.com`);
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setBrandName(e.target.value);
  };

  const handleInputBlur = () => {
    if (brandName.trim() !== '') {
      setDomainVariations(generateDomainVariations(brandName));
      setShowSecondHalf(false);
    }
  };

  const toggleDisplay = () => {
    setShowSecondHalf((prev) => !prev);
  };

  const displayedDomains = showSecondHalf
    ? domainVariations.slice(10, 20)
    : domainVariations.slice(0, 10);

  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='flex items-center justify-between font-medium'>
        Add outbound domains
      </CardHeader>
      <CardContent className='p-0 text-sm'>
        Search and add your ideal outbound domains based on your brand
        <CardFooter className='w-full px-0 mb-0 py-0 mt-2 relative'>
          <Input
            size='sm'
            variant='outline'
            value={brandName}
            placeholder='Brand name'
            onBlur={handleInputBlur}
            onChange={handleInputChange}
          />
          <SearchSm className='absolute right-2 text-gray-500' />
        </CardFooter>
        <div className='mt-2'>
          {displayedDomains.map((domain, index) => (
            <div
              key={index}
              onMouseLeave={() => setIsHovered(null)}
              onMouseEnter={() => setIsHovered(index)}
              className='flex items-center justify-between py-1'
            >
              <span className='text-sm'>{domain}</span>
              {isHovered === index && (
                <IconButton
                  size='xxs'
                  variant='ghost'
                  aria-label='add to cart'
                  icon={<ShoppingCartAdd className='text-primary-700' />}
                />
              )}
            </div>
          ))}
        </div>
        {domainVariations.length > 10 && (
          <Button
            size='xxs'
            variant='ghost'
            colorScheme='primary'
            onClick={toggleDisplay}
            leftIcon={<RefreshCw01 />}
          >
            Suggest new
          </Button>
        )}
      </CardContent>
    </Card>
  );
};

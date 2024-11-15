import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const CheckoutCard = () => {
  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='flex items-center justify-between font-medium text-sm'>
        <span>Base Boundle</span>
        <span>$199.99</span>
      </CardHeader>
      <CardContent>
        <span>1 of 5 domains</span>
      </CardContent>
    </Card>
  );
};

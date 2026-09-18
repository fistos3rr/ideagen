import '@/app/ui/global.css';
import { inter } from '@/app/ui/fonts';
import Header from '@/app/ui/header';
import { headers } from 'next/headers';
import { Suspense } from 'react';


export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const authorized = (await headers()).get('x-authorized') === 'true';

  return (
    <html lang="en">
      <body className={`${inter.className} antialiased`}>
        <Suspense fallback={<div>Loading...</div>}>
          <Header authorized={authorized} />      
          <main>{children}</main>
        <Suspense />
      </body>
    </html>
  );
}

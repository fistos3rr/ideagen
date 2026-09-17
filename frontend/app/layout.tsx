import '@/app/ui/global.css';
import { inter } from '@/app/ui/fonts';
import Header from '@/app/ui/header';
import { headers } from 'next/headers';


export default function RootLayout({ children }: { children: React.ReactNode }) {
  const authorized = (await headers()).get('x-authorized') === 'true';

  return (
    <html lang="en">
      <body className={`${inter.className} antialiased`}>
        <Header authorized={authorized} />      
        <main>{children}</main>
      </body>
    </html>
  );
}

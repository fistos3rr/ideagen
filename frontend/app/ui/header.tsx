import Link from 'next/link';
import { IsAuthorized } from '@/lib/api/client';

export default async function Header() {
  const authorized = await IsAuthorized();

  return (
    <header className="header">
      <nav>
        <Link href="/">Main</Link>
        {authorized ? (
          <Link href="/me">My profile</Link>
        ):(
          <Link href="/login">Login</Link>
        )}
      </nav>
    </header>
  );
}

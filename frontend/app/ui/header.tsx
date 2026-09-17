import Link from 'next/link';
import { isAuthorized } from '@/lib/api/auth';

export default async function Header() {
  const authorized = await isAuthorized();

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

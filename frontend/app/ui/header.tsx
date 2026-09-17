import Link from 'next/link';
import LogoutButton from '@/app/components/LogoutButton'

export default function Header({ authorized }: { authorized: boolean }) {
  return (
    <header className="header">
      <nav>
        <Link href="/">Main</Link>
        {authorized ? (
          <>
            <Link href="/me">My profile</Link>
            <LogoutButton/>
          </>
        ):(
          <Link href="/login">Login</Link>
        )}
      </nav>
    </header>
  );
}

import Link from 'next/link';
import LogoutButton from '@/app/components/LogoutButton'
import { Button } from './button'

export default function Header({ authorized }: { authorized: boolean }) {
  return (
    <header className="header">
      <nav>
        <Link href="/">Main</Link>
        {authorized ? (
          <>
            <Link href="/me">
              <Button>My profile</Button>
            </Link>
            <LogoutButton/>
          </>
        ):(
          <Link href="/login">
            <Button>Login</Button>
          </Link>
        )}
      </nav>
    </header>
  );
}

'use client'

import { Button } from '@/app/ui/button';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function LogoutButton() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const handleLogout = async () => {
    setLoading(true);
    try {
      await fetch('/api/auth/logout', { method: 'POST' });
    } finally {
      setLoading(false);
      router.refresh();
    }
  };

  return (
    <Button onClick={handleLogout} disabled={loading}>
      Logout
    </Button>
  )
}

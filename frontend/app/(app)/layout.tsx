import { requireSession } from '@/lib/auth/session';

export default async function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  await requireSession(); // редирект на /login, если нет сессии
  return <>{children}</>;
}
